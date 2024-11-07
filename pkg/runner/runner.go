package runner

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/abklabs/svmkit/pkg/ssh"
	"github.com/abklabs/svmkit/pkg/utils"
)

// CommandInterface is an interface for representing a script command to executed by the runner.
type Command interface {
	Check() error
	Env() *utils.EnvBuilder
	Script() string
}

// NewRunner initializes a new Runner instance with the given SSH connection.
func NewRunner(conn ssh.Connection, cmd Command) *Runner {
	return &Runner{connection: conn, command: cmd}
}

// Runner represents the setup configuration for a machine.
type Runner struct {
	connection ssh.Connection
	command    Command
}

// Run executes the given setup script on the remote machine.
func (r *Runner) Run(ctx context.Context) error {
	// Load the install script
	scriptBuffer := bytes.NewBufferString(r.command.Script())

	// Generate the environment variables
	var envBuffer bytes.Buffer
	for _, value := range r.command.Env().Args() {
		envBuffer.WriteString(value)
		envBuffer.WriteString("\n")
	}

	// Establish SSH connection
	connection, err := r.connection.Dial(ctx)
	if err != nil {
		return fmt.Errorf("failed to establish SSH connection: %w", err)
	}
	defer connection.Close()

	// Ensure the SSH connection is ready with a context deadline of 10 seconds
	readyCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = ssh.Ready(readyCtx, connection)
	if err != nil {
		return fmt.Errorf("SSH connection not ready: %w", err)
	}

	// Create a temporary directory on the remote host
	tempDir := fmt.Sprintf("/tmp/runner-%d", rand.Int())
	_, _, err = ssh.Exec(ctx, connection, fmt.Sprintf("mkdir -p %s", tempDir))
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}

	logPath := fmt.Sprintf("%s.log", tempDir)

	// Upload the script and environment variables to the remote host
	libPath := path.Join(tempDir, "lib.bash")
	runPath := path.Join(tempDir, "run.sh")
	envPath := path.Join(tempDir, "env")
	stepsPath := path.Join(tempDir, "steps.sh")

	uploads := []struct {
		content []byte
		path    string
	}{
		{content: []byte(LibBash), path: libPath},
		{content: []byte(RunScript), path: runPath},
		{content: envBuffer.Bytes(), path: envPath},
		{content: scriptBuffer.Bytes(), path: stepsPath},
	}

	for _, upload := range uploads {
		err = ssh.Upload(ctx, connection, upload.content, upload.path)
		if err != nil {
			return fmt.Errorf("failed to upload %s: %w", upload.path, err)
		}
	}

	// Make the run.sh script executable, change directory to the temporary directory, execute the script,
	// change back to the original directory, and remove the temporary directory.
	commands := fmt.Sprintf(`
		( chmod +x %s &&
		cd %s &&
		./run.sh &&
		cd - &&
		rm -rf %s ) > %s 2>&1 && rm -f %s
	`, runPath, tempDir, tempDir, logPath, logPath)

	_, _, err = ssh.Exec(ctx, connection, commands)
	if err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}

	return nil
}

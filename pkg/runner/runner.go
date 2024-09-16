package runner

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"path"
	"time"

	"github.com/abklabs/svmkit/pkg/ssh"
)

// InstallCommandInterface is an interface for the InstallCommand struct.
type Command interface {
	Env() map[string]string
	Script() string
}

// Runner represents the setup configuration for a machine.
type Runner struct {
	connection  ssh.Connection
	environment map[string]string
	script      string
}

// Machine initializes a new Runner instance with the given SSH connection.
func Machine(conn ssh.Connection) *Runner {
	return &Runner{connection: conn}
}

// Env sets the environment variables for the setup.
func (r *Runner) Env(entries map[string]string) *Runner {
	r.environment = entries
	return r
}

func (r *Runner) Script(script string) *Runner {
	r.script = script
	return r
}

// Run executes the given validator Runner on the remote machine.
func (r *Runner) Run(ctx context.Context) error {
	// Load the install script
	scriptBuffer := bytes.NewBufferString(r.script)

	// Generate the environment variables
	var envBuffer bytes.Buffer
	for key, value := range r.environment {
		envBuffer.WriteString(fmt.Sprintf("%s=\"%s\"\n", key, value))
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
		chmod +x %s &&
		cd %s &&p
		./run.sh &&
		cd - &&
		rm -rf %s
	`, runPath, tempDir, tempDir)

	_, _, err = ssh.Exec(ctx, connection, commands)
	if err != nil {
		return fmt.Errorf("failed to execute script: %w", err)
	}

	return nil
}

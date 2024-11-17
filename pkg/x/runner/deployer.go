package runner

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type DeployerHandler interface {
	// IngestReaders is responsible for keeping the readers drained.
	// After the readers have been closed, it MUST signal completion by
	// closing the provided done channel.
	IngestReaders(done chan<- struct{}, stdout io.Reader, stderr io.Reader) error
	AugmentError(error) error
}

type Deployer struct {
	Payload *Payload
	Client  *ssh.Client
}

func (p *Deployer) Deploy() error {
	sftpClient, err := sftp.NewClient(p.Client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	defer sftpClient.Close()

	for _, f := range p.Payload.Files {
		path := filepath.Join(p.Payload.RootPath, f.Path)

		dir := filepath.Dir(path)

		if err := sftpClient.MkdirAll(dir); err != nil {
			return fmt.Errorf("failed to create remote directory for %s: %w", dir, err)
		}

		remoteFile, err := sftpClient.Create(path)

		if err != nil {
			return fmt.Errorf("failed to create remote file %s: %w", path, err)
		}

		defer remoteFile.Close()

		if err := remoteFile.Chmod(f.Mode); err != nil {
			return fmt.Errorf("couldn't change ownership of file %s: %w", path, err)
		}

		if _, err := io.Copy(remoteFile, f.Reader); err != nil {
			return fmt.Errorf("failed to write to remote file %s: %w", path, err)
		}
	}

	return nil
}

func (p *Deployer) Run(cmdSegs []string, dontCleanup bool, handler DeployerHandler) error {
	// XXX - This looks worse than it is, but we should come up
	// with a nicer way to do this that isn't Sprintfs.
	cmdSegs = append([]string{"ret=0", ";", "(", "set", "-euo", "pipefail", ";", "cd", p.Payload.RootPath, ";"}, cmdSegs...)
	cmdSegs = append(cmdSegs, []string{")", "||", "ret=$?", ";", "rm", "-rf", p.Payload.RootPath, ";", "exit", "$ret"}...)

	cmd := strings.Join(cmdSegs, " ")

	execSession, err := p.Client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer execSession.Close()

	stdoutPipe, err := execSession.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := execSession.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := execSession.Start(cmd); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	done := make(chan struct{})

	if err := handler.IngestReaders(done, stdoutPipe, stderrPipe); err != nil {
		return fmt.Errorf("couldn't bind command stream handlers: %w", err)
	}

	<-done

	if err := execSession.Wait(); err != nil {
		err = handler.AugmentError(err)
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

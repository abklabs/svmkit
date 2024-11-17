package runner

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type DeployerHandler func(stdout io.Reader, stderr io.Reader) error

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
	cmdSegs = append([]string{"cd", p.Payload.RootPath, "&&"}, cmdSegs...)
	cmdSegs = append(cmdSegs, ";", "rm", "-rf", p.Payload.RootPath)

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

	if err := handler(stdoutPipe, stderrPipe); err != nil {
		return fmt.Errorf("couldn't bind command stream handlers: %w", err)
	}

	if err := execSession.Start(cmd); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	if err := execSession.Wait(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

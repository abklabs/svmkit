package deployer

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/abklabs/svmkit/pkg/runner/payload"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var runWrapperTemplate = template.Must(template.New("runWrapper").Parse(`ret=0 ; ( set -euo pipefail ; cd {{ .RootPath }} ; {{ .Cmd }} ; ) || ret=$? ; {{ if not .KeepPayload }} rm -rf {{ .RootPath }} ; {{ end }} exit $ret`))

type DeployerHandler interface {
	// IngestReaders is responsible for keeping the readers drained.
	// After the readers have been closed, it MUST signal completion by
	// closing the provided done channel.
	IngestReaders(done chan<- struct{}, stdout io.Reader, stderr io.Reader) error
	AugmentError(error) error
}

type SSH struct {
	Payload     *payload.Payload
	Client      *ssh.Client
	KeepPayload bool
}

func (p *SSH) Deploy(statusCallback ProgressStatusCallback) (err error) {
	sftpClient, err := sftp.NewClient(p.Client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}

	defer func() {
		err = errors.Join(err, sftpClient.Close())
	}()

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

		defer func() {
			err = errors.Join(err, remoteFile.Close())
		}()

		if err := remoteFile.Chmod(f.Mode); err != nil {
			return fmt.Errorf("couldn't change ownership of file %s: %w", path, err)
		}

		tracker, err := NewProgressStatus(
			f.Path,
			f.Reader,
			statusCallback)

		if err != nil {
			return fmt.Errorf("couldn't create progress status for %s: %w", path, err)
		}

		if _, err := io.Copy(remoteFile, tracker); err != nil {
			return fmt.Errorf("failed to write to remote file %s: %w", path, err)
		}
	}
	return nil
}

func (p *SSH) Run(cmdSegs []string, handler DeployerHandler) (err error) {
	runWrapper := &strings.Builder{}

	err = runWrapperTemplate.Execute(runWrapper, struct {
		*payload.Payload
		KeepPayload bool
		Cmd         string
	}{
		p.Payload,
		p.KeepPayload,
		strings.Join(cmdSegs, " "),
	})

	if err != nil {
		return fmt.Errorf("couldn't format the deployer's run wrapper: %w", err)
	}

	execSession, err := p.Client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}

	defer func() {
		if closeErr := execSession.Close(); closeErr != io.EOF {
			err = errors.Join(err, closeErr)
		}
	}()

	stdoutPipe, err := execSession.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := execSession.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := execSession.Start(runWrapper.String()); err != nil {
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

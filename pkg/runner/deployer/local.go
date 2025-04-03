package deployer

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner/payload"
)

type Local struct {
	Payload     *payload.Payload
	KeepPayload bool
}

func (p *Local) Deploy() (err error) {
	for _, f := range p.Payload.Files {
		path := filepath.Join(p.Payload.RootPath, f.Path)
		dir := filepath.Dir(path)

		err = os.MkdirAll(dir, 0755)

		if err != nil {
			return err
		}

		file, err := os.Create(path)

		if err != nil {
			return nil
		}

		defer func() {
			err = errors.Join(err, file.Close())
		}()

		err = file.Chmod(f.Mode)

		if err != nil {
			return err
		}

		// Write contents
		if _, err := io.Copy(file, f.Reader); err != nil {
			return fmt.Errorf("failed to write to local file %s: %w", path, err)
		}
	}

	return nil
}

func (p *Local) Run(cmdSegs []string, handler DeployerHandler) error {
	runWrapper := &strings.Builder{}

	err := runWrapperTemplate.Execute(runWrapper, struct {
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

	cmd := exec.Command("bash", "-c", runWrapper.String())

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	done := make(chan struct{})

	if err := handler.IngestReaders(done, stdoutPipe, stderrPipe); err != nil {
		return fmt.Errorf("couldn't bind command stream handlers: %w", err)
	}

	<-done

	if err := cmd.Wait(); err != nil {
		err = handler.AugmentError(err)
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}

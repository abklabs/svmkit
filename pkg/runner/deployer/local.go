package deployer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/abklabs/svmkit/pkg/runner/payload"
)

type Local struct {
	Payload     *payload.Payload
	KeepPayload bool
}

func (p *Local) Deploy() error {
	for _, f := range p.Payload.Files {
		path := filepath.Join(p.Payload.RootPath, f.Path)
		dir := filepath.Dir(path)

		err := os.MkdirAll(dir, 0755)

		if err != nil {
			return err
		}

		file, err := os.Create(path)

		if err != nil {
			return nil
		}

		defer file.Close()

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

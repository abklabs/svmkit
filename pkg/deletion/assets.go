package deletion

import (
	"embed"

	"github.com/abklabs/svmkit/pkg/runner"
)

//go:embed assets
var assets embed.FS

const (
	assetsLib = "assets/lib.sh"
)

func AddToPayload(p *runner.Payload) error {
	file, err := assets.Open(assetsLib)
	if err != nil {
		return err
	}

	p.AddReader("deletion-lib.sh", file)

	return nil
}

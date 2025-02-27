package watchtower

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/watchtower_default_paths.toml
var defaultWatchtowerPathsToml []byte

type WatchtowerPaths struct {
	paths.Paths

	WatchtowerRunBinPath *string `pulumi:"watchtowerRunBinPath,optional" toml:"watchtowerRunBinPath,omitempty"`
}

func NewDefaultWatchtowerPaths() (*WatchtowerPaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &WatchtowerPaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultWatchtowerPathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
}

func (f *WatchtowerPaths) Merge(other *WatchtowerPaths) error {
	if other == nil {
		return nil
	}
	if err := f.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(f, other, mergo.WithOverride)
}

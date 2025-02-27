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

	RunBinPath *string `pulumi:"runBinPath,optional" toml:"runBinPath"`
}

func NewDefaultWatchtowerPaths(base *paths.Paths) (*WatchtowerPaths, error) {
	var b *paths.Paths
	if base == nil {
		def, err := paths.NewDefaultPaths()
		if err != nil {
			return nil, err
		}
		b = def
	} else {
		b = base
	}

	wt := &WatchtowerPaths{
		Paths: *b,
	}

	if err := toml.Unmarshal(defaultWatchtowerPathsToml, wt); err != nil {
		return nil, err
	}
	return wt, nil
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

func (f *WatchtowerPaths) Check() error {
	if err := f.Paths.Check(); err != nil {
		return err
	}
	return paths.CheckPointersNotNil(f)
}

package explorer

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/explorer_default_paths.toml
var defaultExplorerPathsToml []byte

type ExplorerPaths struct {
	paths.Paths

	ExplorerInstallPath *string `pulumi:"explorerInstallPath,optional" toml:"explorerInstallPath,omitempty"`
	ExplorerRunBinPath  *string `pulumi:"explorerRunBinPath,optional" toml:"explorerRunBinPath,omitempty"`
}

func NewDefaultExplorerPaths() (*ExplorerPaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &ExplorerPaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultExplorerPathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
}

func (f *ExplorerPaths) Merge(other *ExplorerPaths) error {
	if other == nil {
		return nil
	}
	if err := f.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(f, other, mergo.WithOverride)
}

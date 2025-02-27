package agave

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/agave_default_paths.toml
var defaultAgavePathsToml []byte

type AgavePaths struct {
	paths.Paths

	AgaveRunBinPath   *string `pulumi:"agaveRunBinPath,optional" toml:"agaveRunBinPath,omitempty"`
	AgaveStopBinPath  *string `pulumi:"agaveStopBinPath,optional" toml:"agaveStopBinPath,omitempty"`
	AgaveCheckBinPath *string `pulumi:"agaveCheckBinPath,optional" toml:"agaveCheckBinPath,omitempty"`
}

func NewDefaultAgavePaths() (*AgavePaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &AgavePaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultAgavePathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
}

func (f *AgavePaths) Merge(other *AgavePaths) error {
	if other == nil {
		return nil
	}
	if err := f.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(f, other, mergo.WithOverride)
}

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

	RunBinPath   *string `pulumi:"runBinPath,optional" toml:"runBinPath"`
	StopBinPath  *string `pulumi:"stopBinPath,optional" toml:"stopBinPath"`
	CheckBinPath *string `pulumi:"checkBinPath,optional" toml:"checkBinPath"`
}

func NewDefaultAgavePaths(base *paths.Paths) (*AgavePaths, error) {
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

	agavePaths := &AgavePaths{
		Paths: *b,
	}

	if err := toml.Unmarshal(defaultAgavePathsToml, agavePaths); err != nil {
		return nil, err
	}
	return agavePaths, nil
}

func (p *AgavePaths) Merge(other *AgavePaths) error {
	if other == nil {
		return nil
	}
	if err := p.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(p, other, mergo.WithOverride)
}

func (p *AgavePaths) MergeFlags(f *AgaveFlags) error {
	if p == nil {
		return nil
	}

	if f.Log != nil {
		p.Paths.LogPath = f.Log
	}

	return nil
}

func (p *AgavePaths) Check() error {
	if err := p.Paths.Check(); err != nil {
		return err
	}
	return paths.CheckPointersNotNil(p)
}

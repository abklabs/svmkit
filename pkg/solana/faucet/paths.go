package faucet

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/faucet_default_paths.toml
var defaultFaucetPathsToml []byte

type FaucetPaths struct {
	paths.Paths

	KeypairPath *string `pulumi:"keypairPath,optional" toml:"keypairPath"`
	RunBinPath  *string `pulumi:"runBinPath,optional" toml:"runBinPath"`
}

func NewDefaultFaucetPaths(base *paths.Paths) (*FaucetPaths, error) {
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

	fp := &FaucetPaths{
		Paths: *b,
	}

	if err := toml.Unmarshal(defaultFaucetPathsToml, fp); err != nil {
		return nil, err
	}
	return fp, nil
}

func (f *FaucetPaths) Merge(other *FaucetPaths) error {
	if other == nil {
		return nil
	}
	if err := f.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(f, other, mergo.WithOverride)
}

func (f *FaucetPaths) Check() error {
	if err := f.Paths.Check(); err != nil {
		return err
	}
	return paths.CheckPointersNotNil(f)
}

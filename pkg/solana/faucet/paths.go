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

	FaucetKeypairPath *string `pulumi:"faucetKeypairPath,optional" toml:"faucetKeypairPath,omitempty"`
	FaucetRunBinPath  *string `pulumi:"faucetRunBinPath,optional" toml:"faucetRunBinPath,omitempty"`
}

func NewDefaultFaucetPaths() (*FaucetPaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &FaucetPaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultFaucetPathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
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

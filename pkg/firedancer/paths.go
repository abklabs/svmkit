package firedancer

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/firedancer_default_paths.toml
var defaultFiredancerPathsToml []byte

type FiredancerPaths struct {
	paths.Paths

	FiredancerConfigPath *string `pulumi:"firedancerConfigPath,optional" toml:"firedancerConfigPath,omitempty"`
}

func NewDefaultFiredancerPaths() (*FiredancerPaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &FiredancerPaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultFiredancerPathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
}

func (f *FiredancerPaths) Merge(other *FiredancerPaths) error {
	if other == nil {
		return nil
	}
	if err := f.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(f, other, mergo.WithOverride)
}

func (f *FiredancerPaths) MergeConfig(c *Config) error {
	if c == nil {
		return nil
	}
	if c.Ledger != nil {
		if c.Ledger.Path == nil && f.LedgerPath != nil {
			c.Ledger.Path = f.LedgerPath
		}
		if c.Ledger.AccountsPath == nil && f.AccountsPath != nil {
			c.Ledger.AccountsPath = f.AccountsPath
		}
	}
	if c.Consensus != nil {
		if c.Consensus.IdentityPath == nil && f.ValidatorIdentityKeypairPath != nil {
			c.Consensus.IdentityPath = f.ValidatorIdentityKeypairPath
		}
		if c.Consensus.VoteAccountPath == nil && f.ValidatorVoteAccountKeypairPath != nil {
			c.Consensus.VoteAccountPath = f.ValidatorVoteAccountKeypairPath
		}
	}
	if c.Log != nil {
		if c.Log.Path == nil && f.LogPath != nil {
			c.Log.Path = f.LogPath
		}
	}
	return nil
}

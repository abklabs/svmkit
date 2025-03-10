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

	ConfigPath *string `pulumi:"configPath,optional" toml:"configPath"`
}

func NewDefaultFiredancerPaths(base *paths.Paths) (*FiredancerPaths, error) {
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

	fd := &FiredancerPaths{
		Paths: *b,
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

func (f *FiredancerPaths) Check() error {
	if err := f.Paths.Check(); err != nil {
		return err
	}
	return paths.CheckPointersNotNil(f)
}

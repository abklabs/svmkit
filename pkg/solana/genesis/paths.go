package genesis

import (
	_ "embed"

	"github.com/abklabs/svmkit/pkg/paths"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/genesis_default_paths.toml
var defaultGenesisPathsToml []byte

type GenesisPaths struct {
	paths.Paths

	GenesisPrimordialAccountsPath *string `pulumi:"genesisPrimordialAccountsPath,optional" toml:"genesisPrimordialAccountsPath,omitempty"`
	GenesisValidatorAccountsPath  *string `pulumi:"genesisValidatorAccountsPath,optional" toml:"genesisValidatorAccountsPath,omitempty"`
	GenesisSolanaSplCachePath     *string `pulumi:"genesisSolanaSplCachePath,optional" toml:"genesisSolanaSplCachePath,omitempty"`
}

func NewDefaultGenesisPaths() (*GenesisPaths, error) {
	base, err := paths.NewDefaultPaths()
	if err != nil {
		return nil, err
	}

	fd := &GenesisPaths{
		Paths: *base,
	}

	if err := toml.Unmarshal(defaultGenesisPathsToml, fd); err != nil {
		return nil, err
	}
	return fd, nil
}

func (p *GenesisPaths) Merge(other *GenesisPaths) error {
	if other == nil {
		return nil
	}
	if err := p.Paths.Merge(&other.Paths); err != nil {
		return err
	}
	return mergo.Merge(p, other, mergo.WithOverride)
}

func (p *GenesisPaths) MergeFlags(flags *GenesisFlags) error {
	if flags == nil {
		return nil
	}

	if flags.LedgerPath != nil && *flags.LedgerPath != "" {
		p.Paths.LedgerPath = flags.LedgerPath
	}

	return nil
}

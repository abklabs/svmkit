package paths

import (
	_ "embed"

	"dario.cat/mergo"
	"github.com/BurntSushi/toml"
)

//go:embed defaults/default_paths.toml
var defaultPathsToml []byte

type Paths struct {
	LogPath                         *string `pulumi:"logPath,optional" toml:"logPath" svmkit:"optional"`
	SystemdPath                     *string `pulumi:"systemdPath,optional" toml:"systemdPath"`
	LedgerPath                      *string `pulumi:"ledgerPath,optional" toml:"ledgerPath"`
	AccountsPath                    *string `pulumi:"accountsPath,optional" toml:"accountsPath"`
	ValidatorIdentityKeypairPath    *string `pulumi:"validatorIdentityKeypairPath,optional" toml:"validatorIdentityKeypairPath"`
	ValidatorVoteAccountKeypairPath *string `pulumi:"validatorVoteAccountKeypairPath,optional" toml:"validatorVoteAccountKeypairPath"`
}

func NewDefaultPaths() (*Paths, error) {
	var fl Paths
	if err := toml.Unmarshal(defaultPathsToml, &fl); err != nil {
		return nil, err
	}
	return &fl, nil
}

func (fl *Paths) Merge(other *Paths) error {
	if other == nil {
		return nil
	}
	return mergo.Merge(fl, other, mergo.WithOverride)
}

func (fl *Paths) Check() error {
	return CheckPointersNotNil(fl)
}

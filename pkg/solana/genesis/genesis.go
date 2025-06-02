package genesis

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/abklabs/svmkit/pkg/deletion"
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"gopkg.in/yaml.v3"
)

const (
	primordialAccountPath  = "/home/sol/primordial.yaml"
	primordialDefaultOwner = "11111111111111111111111111111111"
	validatorAccountsPath  = "/home/sol/validator_accounts.yaml"
)

type CreateCommand struct {
	Genesis
}

func (cmd *CreateCommand) Env() *runner.EnvBuilder {
	genesisEnv := runner.NewEnvBuilder()

	b := runner.NewEnvBuilder()

	b.SetArray("GENESIS_FLAGS", cmd.Flags.Args(cmd.Accounts))
	b.SetArray("GENESIS_ENV", genesisEnv.Args())
	b.Set("LEDGER_PATH", cmd.Flags.LedgerPath)

	b.Merge(cmd.RunnerCommand.Env())

	cmd.DeletionPolicy.Create(&cmd.Genesis, b)

	return b
}

func (cmd *CreateCommand) Check() error {
	if cmd.Flags.HashesPerTick != nil {
		value := *cmd.Flags.HashesPerTick
		switch value {
		case "auto", "sleep":
		default:
			if _, err := strconv.Atoi(value); err != nil {
				return fmt.Errorf("invalid value for HashesPerTick: %q; must be 'auto', 'sleep' or a number", value)
			}
		}
	}

	{
		cmd.SetConfigDefaults()

		packages := deb.Package{}.MakePackageGroup("bzip2")
		packages.Add(deb.Package{Version: cmd.Version}.MakePackages("svmkit-solana-genesis", "svmkit-solana-cli", "svmkit-agave-ledger-tool")...)

		if err := cmd.UpdatePackageGroup(packages); err != nil {
			return err
		}
	}

	policy := cmd.GetDeletionPolicy()
	if err := policy.Check(); err != nil {
		return err
	}

	cmd.DeletionPolicy = &policy

	return nil
}

func (g *Genesis) Create() runner.Command {
	return &CreateCommand{
		Genesis: *g,
	}
}

func (cmd *CreateCommand) AddToPayload(p *runner.Payload) error {
	err := cmd.BuildPrimordialYaml(p.NewWriter(runner.PayloadFile{
		Path: "primordial.yaml",
	}))

	if err != nil {
		return err
	}

	if len(cmd.Accounts) > 0 {
		err = cmd.BuildAccountsYaml(p.NewWriter(runner.PayloadFile{
			Path: "validator_accounts.yaml",
		}))

		if err != nil {
			return err
		}
	}

	genesisScript, err := assets.Open(assetsGenesisScript)
	if err != nil {
		return err
	}

	p.AddReader("steps.sh", genesisScript)

	err = cmd.RunnerCommand.AddToPayload(p)

	if err != nil {
		return err
	}

	if err := deletion.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

func (cmd *CreateCommand) BuildAccountsYaml(w io.Writer) (err error) {
	enc := yaml.NewEncoder(w)

	defer func() {
		err = errors.Join(err, enc.Close())
	}()

	output := map[string][]BootstrapAccount{
		"validator_accounts": cmd.Accounts,
	}

	err = enc.Encode(output)

	if err != nil {
		return fmt.Errorf("failed to encode accounts YAML: %w", err)
	}

	return nil
}

func (cmd *CreateCommand) BuildPrimordialYaml(w io.Writer) (err error) {
	enc := yaml.NewEncoder(w)

	defer func() {
		err = errors.Join(err, enc.Close())
	}()

	resultMap := make(map[string]PrimordialAccount)
	for _, acc := range cmd.Primordial {
		if acc.Owner == "" {
			acc.Owner = primordialDefaultOwner
		}
		resultMap[acc.Pubkey] = acc
	}

	if err = enc.Encode(resultMap); err != nil {
		return fmt.Errorf("failed to encode primordial YAML: %w", err)
	}

	return nil
}

type DeleteCommand struct {
	Genesis
}

// AddToPayload implements runner.Command.
// Subtle: this method shadows the method (Genesis).AddToPayload of DeleteCommand.Genesis.
func (d *DeleteCommand) AddToPayload(p *runner.Payload) error {
	uninstallScript, err := assets.Open(assetsUninstallScript)
	if err != nil {
		return err
	}

	p.AddReader("steps.sh", uninstallScript)

	if err := deletion.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

// Check implements runner.Command.
func (d *DeleteCommand) Check() error {
	d.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup()

	if err := d.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	policy := d.GetDeletionPolicy()
	if err := policy.Check(); err != nil {
		return err
	}

	d.DeletionPolicy = &policy

	return nil
}

// Config implements runner.Command.
// Subtle: this method shadows the method (Genesis).Config of DeleteCommand.Genesis.
func (d *DeleteCommand) Config() *runner.Config {
	return d.RunnerConfig
}

// Env implements runner.Command.
// Subtle: this method shadows the method (Genesis).Env of DeleteCommand.Genesis.
func (d *DeleteCommand) Env() *runner.EnvBuilder {
	b := runner.NewEnvBuilder()

	b.Merge(d.RunnerCommand.Env())

	d.DeletionPolicy.Delete(&d.Genesis, b)

	return b
}

func (g *Genesis) Delete() runner.Command {
	return &DeleteCommand{
		Genesis: *g,
	}
}

// maps to --validator-accounts-file
type BootstrapAccount struct {
	IdentityPubkey  string `pulumi:"identityPubkey" yaml:"identity_account"`
	VotePubkey      string `pulumi:"votePubkey" yaml:"vote_account"`
	StakePubkey     string `pulumi:"stakePubkey" yaml:"stake_account"`
	BalanceLamports *int   `pulumi:"balanceLamports,optional" yaml:"balance_lamports"`
	StakeLamports   *int   `pulumi:"stakeLamports,optional" yaml:"stake_lamports"`
}

// maps to --bootstrap-validator
type BootstrapValidator struct {
	IdentityPubkey string `pulumi:"identityPubkey" yaml:"identity_account"`
	VotePubkey     string `pulumi:"votePubkey" yaml:"vote_account"`
	StakePubkey    string `pulumi:"stakePubkey" yaml:"stake_account"`
}

type PrimordialAccount struct {
	Pubkey     string `pulumi:"pubkey" yaml:"-"`
	Lamports   int64  `pulumi:"lamports" yaml:"balance"`
	Owner      string `pulumi:"owner,optional" yaml:"owner"`
	Executable bool   `pulumi:"executable,optional" yaml:"executable"`
	Data       string `pulumi:"data,optional" yaml:"data"`
}

type Genesis struct {
	runner.RunnerCommand

	Flags          GenesisFlags        `pulumi:"flags"`
	Primordial     []PrimordialAccount `pulumi:"primordial"`
	Accounts       []BootstrapAccount  `pulumi:"accounts,optional"`
	Version        *string             `pulumi:"version,optional"`
	DeletionPolicy *deletion.Policy    `pulumi:"deletionPolicy,optional"`
}

func (g *Genesis) GetDeletionPolicy() deletion.Policy {
	if g.DeletionPolicy == nil {
		return deletion.PolicyKeep
	} else {
		return *g.DeletionPolicy
	}
}

func (g *Genesis) ManagedFiles() []string {
	return []string{g.Flags.LedgerPath}
}

type GenesisFlags struct {
	LedgerPath string `pulumi:"ledgerPath"`

	BootstrapValidators             []BootstrapValidator `pulumi:"bootstrapValidators"`
	BootstrapStakeAuthorizedPubkey  *string              `pulumi:"bootstrapStakeAuthorizedPubkey,optional"`
	BootstrapValidatorLamports      *int                 `pulumi:"bootstrapValidatorLamports,optional"`
	BootstrapValidatorStakeLamports *int                 `pulumi:"bootstrapValidatorStakeLamports,optional"`
	ClusterType                     *string              `pulumi:"clusterType,optional"`
	CreationTime                    *string              `pulumi:"creationTime,optional"`
	DeactivateFeatures              *[]string            `pulumi:"deactivateFeatures,optional"`
	EnableWarmupEpochs              *bool                `pulumi:"enableWarmupEpochs,optional"`
	FaucetPubkey                    *string              `pulumi:"faucetPubkey,optional"`
	FaucetLamports                  *int                 `pulumi:"faucetLamports,optional"`
	FeeBurnPercentage               *int                 `pulumi:"feeBurnPercentage,optional"`
	HashesPerTick                   *string              `pulumi:"hashesPerTick,optional"` // can be "auto", "sleep", or a number
	Inflation                       *string              `pulumi:"inflation,optional"`
	LamportsPerByteYear             *int                 `pulumi:"lamportsPerByteYear,optional"`
	MaxGenesisArchiveUnpackedSize   *int                 `pulumi:"maxGenesisArchiveUnpackedSize,optional"`
	RentBurnPercentage              *int                 `pulumi:"rentBurnPercentage,optional"`
	RentExemptionThreshold          *int                 `pulumi:"rentExemptionThreshold,optional"`
	SlotsPerEpoch                   *int                 `pulumi:"slotsPerEpoch,optional"`
	TargetLamportsPerSignature      *int                 `pulumi:"targetLamportsPerSignature,optional"`
	TargetSignaturesPerSlot         *int                 `pulumi:"targetSignaturesPerSlot,optional"`
	TargetTickDuration              *int                 `pulumi:"targetTickDuration,optional"`
	TicksPerSlot                    *int                 `pulumi:"ticksPerSlot,optional"`
	Url                             *string              `pulumi:"url,optional"`
	VoteCommissionPercentage        *int                 `pulumi:"voteCommissionPercentage,optional"`
	ExtraFlags                      *[]string            `pulumi:"extraFlags,optional"`
}

func (f GenesisFlags) Args(accounts []BootstrapAccount) []string {
	b := runner.FlagBuilder{}

	// Note: --upgradeable-program, --bpf-program are hard-coded in the install
	// script and should not be included here.

	// Required flags
	b.Append("primordial-accounts-file", primordialAccountPath)

	for _, validator := range f.BootstrapValidators {
		b.AppendRaw("--bootstrap-validator", validator.IdentityPubkey, validator.VotePubkey, validator.StakePubkey)
	}

	if len(accounts) > 0 {
		b.Append("validator-accounts-file", validatorAccountsPath)
	}

	b.Append("ledger", f.LedgerPath)

	// Optional flags
	b.AppendP("bootstrap-stake-authorized-pubkey", f.BootstrapStakeAuthorizedPubkey)
	b.AppendIntP("bootstrap-validator-lamports", f.BootstrapValidatorLamports)
	b.AppendIntP("bootstrap-validator-stake-lamports", f.BootstrapValidatorStakeLamports)
	b.AppendP("cluster-type", f.ClusterType)
	b.AppendP("creation-time", f.CreationTime)
	b.AppendArrayP("deactivate-feature", f.DeactivateFeatures)
	b.AppendBoolP("enable-warmup-epochs", f.EnableWarmupEpochs)
	b.AppendP("faucet-pubkey", f.FaucetPubkey)
	b.AppendIntP("faucet-lamports", f.FaucetLamports)
	b.AppendIntP("fee-burn-percentage", f.FeeBurnPercentage)
	b.AppendP("hashes-per-tick", f.HashesPerTick) // This can be "auto", "sleep" or a number
	b.AppendP("inflation", f.Inflation)
	b.AppendIntP("lamports-per-byte-year", f.LamportsPerByteYear)
	b.AppendIntP("max-genesis-archive-unpacked-size", f.MaxGenesisArchiveUnpackedSize)
	b.AppendIntP("rent-burn-percentage", f.RentBurnPercentage)
	b.AppendIntP("rent-exemption-threshold", f.RentExemptionThreshold)
	b.AppendIntP("slots-per-epoch", f.SlotsPerEpoch)
	b.AppendIntP("target-lamports-per-signature", f.TargetLamportsPerSignature)
	b.AppendIntP("target-signatures-per-slot", f.TargetSignaturesPerSlot)
	b.AppendIntP("target-tick-duration", f.TargetTickDuration)
	b.AppendIntP("ticks-per-slot", f.TicksPerSlot)
	b.AppendP("url", f.Url)
	b.AppendIntP("vote-commission-percentage", f.VoteCommissionPercentage)

	if f.ExtraFlags != nil {
		b.AppendRaw(*f.ExtraFlags...)
	}

	return b.Args()
}

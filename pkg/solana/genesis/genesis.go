package genesis

import (
	"fmt"
	"io"
	"strconv"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"gopkg.in/yaml.v3"
)

const (
	primordialAccountPath  = "/home/sol/primordial.yaml"
	primordialDefaultOwner = "11111111111111111111111111111111"
	validatorAccountsPath  = "/home/sol/validator_accounts.yaml"

	defaultClusterType                = "development"
	defaultFaucetLamports             = 1000
	defaultTargetLamportsPerSignature = 0
	defaultInflation                  = "none"
	defaultLamportsPerByteYear        = 1
	defaultSlotsPerEpoch              = 150
)

type CreateCommand struct {
	Genesis
}

func (cmd *CreateCommand) Env() *runner.EnvBuilder {
	genesisEnv := runner.NewEnvBuilder()

	b := runner.NewEnvBuilder()

	b.SetArray("GENESIS_FLAGS", cmd.Flags.Args(cmd.Accounts, cmd.Paths))
	b.SetArray("GENESIS_ENV", genesisEnv.Args())

	b.Merge(cmd.RunnerCommand.Env())

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
		cmd.RunnerCommand.SetConfigDefaults()

		packages := deb.Package{}.MakePackageGroup("bzip2")
		packages.Add(deb.Package{Version: cmd.Version}.MakePackages("svmkit-solana-genesis", "svmkit-solana-cli", "svmkit-agave-ledger-tool")...)

		if err := cmd.RunnerCommand.UpdatePackageGroup(packages); err != nil {
			return err
		}
	}

	if err := cmd.Paths.MergeFlags(&cmd.Flags); err != nil {
		return err
	}

	if err := cmd.Paths.Check(); err != nil {
		return err
	}

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

	if err := p.AddTemplate("steps.sh", genesisScriptTmpl, cmd); err != nil {
		return err
	}

	err = cmd.RunnerCommand.AddToPayload(p)

	if err != nil {
		return err
	}

	return nil
}

func (cmd *CreateCommand) BuildAccountsYaml(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()

	output := map[string][]BootstrapAccount{
		"validator_accounts": cmd.Accounts,
	}

	if err := enc.Encode(output); err != nil {
		return fmt.Errorf("failed to encode accounts YAML: %w", err)
	}
	return nil
}

func (cmd *CreateCommand) BuildPrimordialYaml(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()

	resultMap := make(map[string]PrimordialAccount)
	for _, acc := range cmd.Primordial {
		if acc.Owner == "" {
			acc.Owner = primordialDefaultOwner
		}
		resultMap[acc.Pubkey] = acc
	}

	if err := enc.Encode(resultMap); err != nil {
		return fmt.Errorf("failed to encode primordial YAML: %w", err)
	}

	return nil
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

	Flags      GenesisFlags        `pulumi:"flags"`
	Paths      GenesisPaths        `pulumi:"paths"`
	Primordial []PrimordialAccount `pulumi:"primordial"`
	Accounts   []BootstrapAccount  `pulumi:"accounts,optional"`
	Version    *string             `pulumi:"version,optional"`
}

type GenesisFlags struct {
	LedgerPath *string `pulumi:"ledgerPath,optional"`

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

func (f GenesisFlags) Args(accounts []BootstrapAccount, paths GenesisPaths) []string {
	b := runner.FlagBuilder{}

	// Note: --upgradeable-program, --bpf-program are hard-coded in the install
	// script and should not be included here.

	// Required flags
	b.AppendP("primordial-accounts-file", paths.PrimordialAccountsPath)

	for _, validator := range f.BootstrapValidators {
		b.AppendRaw("--bootstrap-validator", validator.IdentityPubkey, validator.VotePubkey, validator.StakePubkey)
	}

	if len(accounts) > 0 {
		b.AppendP("validator-accounts-file", paths.ValidatorAccountsPath)
	}

	b.AppendP("ledger", paths.LedgerPath)

	// Optional flags
	b.AppendP("bootstrap-stake-authorized-pubkey", f.BootstrapStakeAuthorizedPubkey)
	b.AppendIntP("bootstrap-validator-lamports", f.BootstrapValidatorLamports)
	b.AppendIntP("bootstrap-validator-stake-lamports", f.BootstrapValidatorStakeLamports)

	if f.ClusterType != nil {
		b.AppendP("cluster-type", f.ClusterType)
	} else {
		value := defaultClusterType
		b.AppendP("cluster-type", &value)
	}

	b.AppendP("creation-time", f.CreationTime)
	b.AppendArrayP("deactivate-feature", f.DeactivateFeatures)
	b.AppendBoolP("enable-warmup-epochs", f.EnableWarmupEpochs)
	b.AppendP("faucet-pubkey", f.FaucetPubkey)

	if f.FaucetLamports != nil {
		b.AppendIntP("faucet-lamports", f.FaucetLamports)
	} else {
		value := defaultFaucetLamports
		b.AppendIntP("faucet-lamports", &value)
	}

	b.AppendIntP("fee-burn-percentage", f.FeeBurnPercentage)
	b.AppendP("hashes-per-tick", f.HashesPerTick) // This can be "auto", "sleep" or a number

	if f.Inflation != nil {
		b.AppendP("inflation", f.Inflation)
	} else {
		value := defaultInflation
		b.AppendP("inflation", &value)
	}

	if f.LamportsPerByteYear != nil {
		b.AppendIntP("lamports-per-byte-year", f.LamportsPerByteYear)
	} else {
		value := defaultLamportsPerByteYear
		b.AppendIntP("lamports-per-byte-year", &value)
	}

	b.AppendIntP("max-genesis-archive-unpacked-size", f.MaxGenesisArchiveUnpackedSize)
	b.AppendIntP("rent-burn-percentage", f.RentBurnPercentage)
	b.AppendIntP("rent-exemption-threshold", f.RentExemptionThreshold)

	if f.SlotsPerEpoch != nil {
		b.AppendIntP("slots-per-epoch", f.SlotsPerEpoch)
	} else {
		value := defaultSlotsPerEpoch
		b.AppendIntP("slots-per-epoch", &value)
	}

	if f.TargetLamportsPerSignature != nil {
		b.AppendIntP("target-lamports-per-signature", f.TargetLamportsPerSignature)
	} else {
		value := defaultTargetLamportsPerSignature
		b.AppendIntP("target-lamports-per-signature", &value)
	}

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

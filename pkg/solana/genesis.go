package solana

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/abklabs/svmkit/pkg/deb"
	"github.com/abklabs/svmkit/pkg/runner"
)

const (
	primordialAccountPath = "/home/sol/primordial.yaml"

	defaultClusterType                = "development"
	defaultFaucetLamports             = 1000
	defaultTargetLamportsPerSignature = 0
	defaultInflation                  = "none"
	defaultLamportsPerByteYear        = 1
	defaultSlotPerEpoch               = 150
)

type CreateCommand struct {
	Genesis
}

func (cmd *CreateCommand) Env() *runner.EnvBuilder {
	genesisEnv := runner.NewEnvBuilder()

	b := runner.NewEnvBuilder()

	b.SetArray("GENESIS_FLAGS", cmd.Flags.Args())
	b.SetArray("GENESIS_ENV", genesisEnv.Args())

	// Primordial accounts as environment variables
	var primordialPubkeys, primordialLamports string
	if cmd.Primordial != nil {
		var pubkeys, lamports []string
		for _, entry := range cmd.Primordial {
			pubkeys = append(pubkeys, entry.Pubkey)
			lamports = append(lamports, entry.Lamports)
		}
		primordialPubkeys = strings.Join(pubkeys, ",")
		primordialLamports = strings.Join(lamports, ",")
	}

	b.Set("LEDGER_PATH", cmd.Flags.LedgerPath)
	b.Set("PRIMORDIAL_PUBKEYS", primordialPubkeys)
	b.Set("PRIMORDIAL_LAMPORTS", primordialLamports)

	{
		packages := deb.Package{}.MakePackageGroup("bzip2")
		packages.Add(deb.Package{Version: cmd.Version}.MakePackages("svmkit-solana-genesis", "svmkit-solana-cli", "svmkit-agave-ledger-tool")...)
		b.SetArray("PACKAGE_LIST", packages.Args())
	}

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

	return nil
}

func (g *Genesis) Create() runner.Command {
	return &CreateCommand{
		Genesis: *g,
	}
}

func (cmd *CreateCommand) AddToPayload(p *runner.Payload) error {
	genesisScript, err := assets.Open(assetsGenesisScript)
	if err != nil {
		return err
	}

	p.AddReader("steps.sh", genesisScript)
	return nil
}

type PrimorialEntry struct {
	Pubkey   string `pulumi:"pubkey"`
	Lamports string `pulumi:"lamports"`
}

type Genesis struct {
	Flags      GenesisFlags     `pulumi:"flags"`
	Primordial []PrimorialEntry `pulumi:"primordial"`
	Version    *string          `pulumi:"version,optional"`
}

type GenesisFlags struct {
	IdentityPubkey string `pulumi:"identityPubkey"`
	LedgerPath     string `pulumi:"ledgerPath"`
	VotePubkey     string `pulumi:"votePubkey"`
	StakePubkey    string `pulumi:"stakePubkey"`

	BootstrapStakeAuthorizedPubkey  *string   `pulumi:"bootstrapStakeAuthorizedPubkey,optional"`
	BootstrapValidatorLamports      *int      `pulumi:"bootstrapValidatorLamports,optional"`
	BootstrapValidatorStakeLamports *int      `pulumi:"bootstrapValidatorStakeLamports,optional"`
	ClusterType                     *string   `pulumi:"clusterType,optional"`
	CreationTime                    *string   `pulumi:"creationTime,optional"`
	DeactivateFeatures              *[]string `pulumi:"deactivateFeatures,optional"`
	EnableWarmupEpochs              *bool     `pulumi:"enableWarmupEpochs,optional"`
	FaucetPubkey                    *string   `pulumi:"faucetPubkey,optional"`
	FaucetLamports                  *int      `pulumi:"faucetLamports,optional"`
	FeeBurnPercentage               *int      `pulumi:"feeBurnPercentage,optional"`
	HashesPerTick                   *string   `pulumi:"hashesPerTick,optional"` // can be "auto", "sleep", or a number
	Inflation                       *string   `pulumi:"inflation,optional"`
	LamportsPerByteYear             *int      `pulumi:"lamportsPerByteYear,optional"`
	MaxGenesisArchiveUnpackedSize   *int      `pulumi:"maxGenesisArchiveUnpackedSize,optional"`
	RentBurnPercentage              *int      `pulumi:"rentBurnPercentage,optional"`
	RentExemptionThreshold          *int      `pulumi:"rentExemptionThreshold,optional"`
	SlotPerEpoch                    *int      `pulumi:"slotPerEpoch,optional"`
	TargetLamportsPerSignature      *int      `pulumi:"targetLamportsPerSignature,optional"`
	TargetSignaturesPerSlot         *int      `pulumi:"targetSignaturesPerSlot,optional"`
	TargetTickDuration              *int      `pulumi:"targetTickDuration,optional"`
	TicksPerSlot                    *int      `pulumi:"ticksPerSlot,optional"`
	Url                             *string   `pulumi:"url,optional"`
	VoteCommissionPercentage        *int      `pulumi:"voteCommissionPercentage,optional"`
	ExtraFlags                      *[]string `pulumi:"extraFlags,optional"`
}

func (f GenesisFlags) Args() []string {
	b := runner.FlagBuilder{}

	// Note: --upgradeable-program, --bpf-program are hard-coded in the install
	// script and should not be included here.

	// Required flags
	b.Append("primordial-accounts-file", primordialAccountPath)
	b.AppendRaw("--bootstrap-validator", f.IdentityPubkey, f.VotePubkey, f.StakePubkey)
	b.Append("ledger", f.LedgerPath)

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

	if f.SlotPerEpoch != nil {
		b.AppendIntP("slots-per-epoch", f.SlotPerEpoch)
	} else {
		value := defaultSlotPerEpoch
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

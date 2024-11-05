package solana

import (
	"strings"

	"github.com/abklabs/svmkit/pkg/genesis"
	"github.com/abklabs/svmkit/pkg/runner"
)

// GenesisFlags represents the configuration flags for the Solana genesis setup.
type GenesisFlags struct {
	LedgerPath                 string  `pulumi:"ledgerPath"`
	IdentityPubkey             string  `pulumi:"identityPubkey"`
	VotePubkey                 string  `pulumi:"votePubkey"`
	StakePubkey                string  `pulumi:"stakePubkey"`
	FaucetPubkey               string  `pulumi:"faucetPubkey"`
	FaucetLamports             *string `pulumi:"faucetLamports,optional"`
	TargetLamportsPerSignature *string `pulumi:"targetLamportsPerSignature,optional"`
	Inflation                  *string `pulumi:"inflation,optional"`
	LamportsPerByteYear        *string `pulumi:"lamportsPerByteYear,optional"`
	SlotPerEpoch               *string `pulumi:"slotPerEpoch,optional"`
	ClusterType                *string `pulumi:"clusterType,optional"`
}

type CreateCommand struct {
	Genesis
}

func (cmd *CreateCommand) Env() map[string]string {
	env := map[string]string{
		"LEDGER_PATH":                   cmd.Flags.LedgerPath,
		"IDENTITY_PUBKEY":               cmd.Flags.IdentityPubkey,
		"VOTE_PUBKEY":                   cmd.Flags.VotePubkey,
		"STAKE_PUBKEY":                  cmd.Flags.StakePubkey,
		"FAUCET_PUBKEY":                 cmd.Flags.FaucetPubkey,
		"FAUCET_LAMPORTS":               "1000",
		"TARGET_LAMPORTS_PER_SIGNATURE": "0",
		"INFLATION":                     "none",
		"LAMPORTS_PER_BYTE_YEAR":        "1",
		"SLOT_PER_EPOCH":                "150",
		"CLUSTER_TYPE":                  "development",
	}

	if cmd.Flags.FaucetLamports != nil {
		env["FAUCET_LAMPORTS"] = *cmd.Flags.FaucetLamports
	}
	if cmd.Flags.TargetLamportsPerSignature != nil {
		env["TARGET_LAMPORTS_PER_SIGNATURE"] = *cmd.Flags.TargetLamportsPerSignature
	}
	if cmd.Flags.Inflation != nil {
		env["INFLATION"] = *cmd.Flags.Inflation
	}
	if cmd.Flags.LamportsPerByteYear != nil {
		env["LAMPORTS_PER_BYTE_YEAR"] = *cmd.Flags.LamportsPerByteYear
	}
	if cmd.Flags.SlotPerEpoch != nil {
		env["SLOT_PER_EPOCH"] = *cmd.Flags.SlotPerEpoch
	}
	if cmd.Flags.ClusterType != nil {
		env["CLUSTER_TYPE"] = *cmd.Flags.ClusterType
	}
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
	env["PRIMORDIAL_PUBKEYS"] = primordialPubkeys
	env["PRIMORDIAL_LAMPORTS"] = primordialLamports

	if cmd.Version != nil {
		env["PACKAGE_VERSION"] = *cmd.Version
	}

	return env
}

func (cmd *CreateCommand) Script() string {
	return GenesisScript
}

type Genesis struct {
	Flags      GenesisFlags
	Primordial []genesis.PrimorialEntry
	Version    genesis.Version
}

func (g *Genesis) Create() runner.Command {
	return &CreateCommand{
		Genesis: *g,
	}
}

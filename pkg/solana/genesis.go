package solana

import (
	"log"
	"os"
	"strings"

	"github.com/abklabs/svmkit/pkg/genesis"
	"github.com/abklabs/svmkit/pkg/runner"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

// centralized logger setup
func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

type CreateCommand struct {
	runner.Command
	Flags      GenesisFlags
	Primordial []genesis.PrimorialEntry
	Version    genesis.Version
}

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

func setEnvFlags(env map[string]string, flags GenesisFlags) {
	flagsMappings := map[string]*string{
		"FAUCET_LAMPORTS":               flags.FaucetLamports,
		"TARGET_LAMPORTS_PER_SIGNATURE": flags.TargetLamportsPerSignature,
		"INFLATION":                     flags.Inflation,
		"LAMPORTS_PER_BYTE_YEAR":        flags.LamportsPerByteYear,
		"SLOT_PER_EPOCH":                flags.SlotPerEpoch,
		"CLUSTER_TYPE":                  flags.ClusterType,
	}

	for key, value := range flagsMappings {
		if value == nil {
			WarningLogger.Printf("Warning: Missing value for environment variable '%s'.", key)
			continue
		}
		env[key] = *value
		InfoLogger.Printf("Set environment variable '%s' to '%s'", key, *value)
	}
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

	setEnvFlags(env, cmd.Flags)

	//considering abstracting this logic if Primodial data would be needed elsewhere
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
	genesis.Genesis
	Flags      GenesisFlags
	Primordial []genesis.PrimorialEntry
	Version    genesis.Version
}

func (g *Genesis) Create() runner.Command {
	return &CreateCommand{
		Flags:      g.Flags,
		Primordial: g.Primordial,
		Version:    g.Version,
	}
}

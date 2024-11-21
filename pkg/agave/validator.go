package agave

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/solana"
	"github.com/abklabs/svmkit/pkg/validator"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const (
	accountsPath = "/home/sol/accounts"
	ledgerPath   = "/home/sol/ledger"
	logPath      = "/home/sol/log"

	identityKeyPairPath    = "/home/sol/validator-keypair.json"
	voteAccountKeyPairPath = "/home/sol/vote-account-keypair.json"
)

type Variant string

const (
	VariantSolana      Variant = "solana"
	VariantAgave       Variant = "agave"
	VariantPowerledger Variant = "powerledger"
	VariantJito        Variant = "jito"
	VariantPyth        Variant = "pyth"
	VariantMantis      Variant = "mantis"
)

func (Variant) Values() []infer.EnumValue[Variant] {
	return []infer.EnumValue[Variant]{
		{
			Name:        string(VariantSolana),
			Value:       VariantSolana,
			Description: "The Solana validator",
		},
		{
			Name:        string(VariantAgave),
			Value:       VariantAgave,
			Description: "The Agave validator",
		},
		{
			Name:        string(VariantPowerledger),
			Value:       VariantPowerledger,
			Description: "The Powerledger validator",
		},
		{
			Name:        string(VariantJito),
			Value:       VariantJito,
			Description: "The Jito validator",
		},
		{
			Name:        string(VariantPyth),
			Value:       VariantPyth,
			Description: "The Pyth validator",
		},
		{
			Name:        string(VariantMantis),
			Value:       VariantMantis,
			Description: "The Mantis validator",
		},
	}
}

type KeyPairs struct {
	Identity    string `pulumi:"identity" provider:"secret"`
	VoteAccount string `pulumi:"voteAccount" provider:"secret"`
}

type Metrics struct {
	URL      string `pulumi:"url"`
	Database string `pulumi:"database"`
	User     string `pulumi:"user"`
	Password string `pulumi:"password"`
}

func (m *Metrics) Check() error {
	if m.URL == "" {
		return fmt.Errorf("metrics URL cannot be empty")
	}

	if m.Database == "" {
		return fmt.Errorf("metrics database cannot be empty")
	}

	if m.User == "" {
		return fmt.Errorf("metrics user cannot be empty")
	}

	return nil
}

// String constructs the Solana metrics configuration string from the separate fields
// and returns it as an environment variable string.
func (m *Metrics) String() string {

	// Note: We allow empty password as it might be a valid case in some scenarios
	configParts := []string{
		fmt.Sprintf("host=%s", m.URL),
		fmt.Sprintf("db=%s", m.Database),
		fmt.Sprintf("u=%s", m.User),
		fmt.Sprintf("p=%s", m.Password),
	}

	return strings.Join(configParts, ",")

}

type InstallCommand struct {
	Agave
}

func (cmd *InstallCommand) Check() error {
	if m := cmd.Metrics; m != nil {
		if err := m.Check(); err != nil {
			return fmt.Errorf("Warning: Invalid metrics URL: %v\n", err)
		}
	}

	return nil
}

func (cmd *InstallCommand) Env() *runner.EnvBuilder {
	validatorEnv := runner.NewEnvBuilder()

	if m := cmd.Metrics; m != nil {
		validatorEnv.Set("SOLANA_METRICS_CONFIG", m.String())
	}

	b := runner.NewEnvBuilder()

	b.SetMap(map[string]string{
		"VALIDATOR_FLAGS": strings.Join(cmd.Flags.ToArgs(), " "),
		"VALIDATOR_ENV":   validatorEnv.String(),
	})

	{
		s := identityKeyPairPath
		conf := solana.CLIConfig{
			KeyPair: &s,
		}

		if senv := cmd.Environment; senv != nil {
			conf.URL = senv.RPCURL
		}

		b.Set("SOLANA_CLI_CONFIG_FLAGS", conf.ToFlags().String())
	}

	b.SetP("VALIDATOR_VERSION", cmd.Version)

	if cmd.Variant != nil {
		b.Set("VALIDATOR_VARIANT", string(*cmd.Variant))
	} else {
		b.Set("VALIDATOR_VARIANT", string(VariantAgave))
	}

	b.SetBoolP("FULL_RPC", cmd.Flags.FullRpcAPI)
	b.Set("RPC_BIND_ADDRESS", cmd.Flags.RpcBindAddress)
	b.SetInt("RPC_PORT", cmd.Flags.RpcPort)

	return b
}

func (cmd *InstallCommand) AddToPayload(p *runner.Payload) error {
	p.AddString("steps.sh", InstallScript)

	p.AddString("validator-keypair.json", cmd.KeyPairs.Identity)
	p.AddString("vote-account-keypair.json", cmd.KeyPairs.VoteAccount)

	return nil
}

type Agave struct {
	Environment *solana.Environment `pulumi:"environment,optional"`
	Version     validator.Version   `pulumi:"version,optional"`
	Variant     *Variant            `pulumi:"variant,optional"`
	KeyPairs    KeyPairs            `pulumi:"keyPairs"`
	Flags       Flags               `pulumi:"flags"`
	Metrics     *Metrics            `pulumi:"metrics,optional"`
}

func (agave *Agave) Install() runner.Command {
	return &InstallCommand{
		Agave: *agave,
	}
}

type Flags struct {
	EntryPoint                   *[]string `pulumi:"entryPoint,optional"`
	KnownValidator               *[]string `pulumi:"knownValidator,optional"`
	UseSnapshotArchivesAtStartup string    `pulumi:"useSnapshotArchivesAtStartup"`
	RpcPort                      int       `pulumi:"rpcPort"`
	PrivateRPC                   bool      `pulumi:"privateRPC"`
	OnlyKnownRPC                 bool      `pulumi:"onlyKnownRPC"`
	DynamicPortRange             string    `pulumi:"dynamicPortRange"`
	GossipHost                   *string   `pulumi:"gossipHost,optional"`
	GossipPort                   int       `pulumi:"gossipPort"`
	RpcBindAddress               string    `pulumi:"rpcBindAddress"`
	WalRecoveryMode              string    `pulumi:"walRecoveryMode"`
	LimitLedgerSize              int       `pulumi:"limitLedgerSize"`
	BlockProductionMethod        string    `pulumi:"blockProductionMethod"`
	TvuReceiveThreads            *int      `pulumi:"tvuReceiveThreads,optional"`
	NoWaitForVoteToStartLeader   bool      `pulumi:"noWaitForVoteToStartLeader"`
	FullSnapshotIntervalSlots    int       `pulumi:"fullSnapshotIntervalSlots"`
	ExpectedGenesisHash          *string   `pulumi:"expectedGenesisHash,optional"`
	FullRpcAPI                   *bool     `pulumi:"fullRpcAPI,optional"`
	NoVoting                     *bool     `pulumi:"noVoting,optional"`
	AllowPrivateAddr             *bool     `pulumi:"allowPrivateAddr,optional"`
	ExtraFlags                   *[]string `pulumi:"extraFlags,optional"`
}

func (f Flags) ToArgs() []string {
	b := runner.FlagBuilder{}

	// Note: These locations are hard coded inside asset-builder.
	b.Append("--identity", identityKeyPairPath)
	b.Append("--vote-account", voteAccountKeyPairPath)

	if f.EntryPoint != nil {
		for _, entrypoint := range *f.EntryPoint {
			b.AppendP("entrypoint", &entrypoint)
		}
	}

	if f.KnownValidator != nil {
		for _, knownValidator := range *f.KnownValidator {
			b.AppendP("known-validator", &knownValidator)
		}
	}

	b.AppendP("expected-genesis-hash", f.ExpectedGenesisHash)

	b.AppendP("use-snapshot-archives-at-startup", &f.UseSnapshotArchivesAtStartup)
	b.AppendIntP("rpc-port", &f.RpcPort)
	b.AppendP("dynamic-port-range", &f.DynamicPortRange)

	b.AppendP("gossip-host", f.GossipHost)

	b.AppendIntP("gossip-port", &f.GossipPort)
	b.AppendP("rpc-bind-address", &f.RpcBindAddress)
	b.AppendP("wal-recovery-mode", &f.WalRecoveryMode)
	b.Append("--log", logPath)
	b.Append("--accounts", accountsPath)
	b.Append("--ledger", ledgerPath)
	b.AppendIntP("limit-ledger-size", &f.LimitLedgerSize)
	b.AppendP("block-production-method", &f.BlockProductionMethod)

	b.AppendIntP("tvu-receive-threads", f.TvuReceiveThreads)

	b.AppendIntP("full-snapshot-interval-slots", &f.FullSnapshotIntervalSlots)
	b.AppendBoolP("no-wait-for-vote-to-start-leader", &f.NoWaitForVoteToStartLeader)
	b.AppendBoolP("only-known-rpc", &f.OnlyKnownRPC)
	b.AppendBoolP("private-rpc", &f.PrivateRPC)

	b.AppendBoolP("full-rpc-api", f.FullRpcAPI)

	b.AppendBoolP("no-voting", f.NoVoting)
	b.AppendBoolP("allow-private-addr", f.AllowPrivateAddr)

	if f.ExtraFlags != nil {
		b.Append(*f.ExtraFlags...)
	}

	return b.ToArgs()
}

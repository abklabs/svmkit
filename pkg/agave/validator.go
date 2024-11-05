package agave

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/validator"
	"github.com/pulumi/pulumi-go-provider/infer"
)

const (
	accountsPath = "/home/sol/accounts"
	ledgerPath   = "/home/sol/ledger"
	logPath      = "/home/sol/log"
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

// ValidatorEnv represents the runtime environment specifically for the validator
type ValidatorEnv struct {
	Metrics *Metrics
}

func (env *ValidatorEnv) ToString() string {
	var envStrings []string

	if env.Metrics != nil {
		metricsEnv, err := env.Metrics.ToEnv()
		if err == nil {
			envStrings = append(envStrings, metricsEnv)
		} else {
			fmt.Printf("Warning: Invalid metrics URL: %v\n", err)
		}
	}

	return strings.Join(envStrings, " ")
}

// ToEnv constructs the Solana metrics configuration string from the separate fields
// and returns it as an environment variable string.
func (m *Metrics) ToEnv() (string, error) {
	if m.URL == "" {
		return "", fmt.Errorf("metrics URL cannot be empty")
	}

	if m.Database == "" {
		return "", fmt.Errorf("metrics database cannot be empty")
	}

	if m.User == "" {
		return "", fmt.Errorf("metrics user cannot be empty")
	}

	// Note: We allow empty password as it might be a valid case in some scenarios
	configParts := []string{
		fmt.Sprintf("host=%s", m.URL),
		fmt.Sprintf("db=%s", m.Database),
		fmt.Sprintf("u=%s", m.User),
		fmt.Sprintf("p=%s", m.Password),
	}

	metricsConfig := strings.Join(configParts, ",")
	// XXX - We should quote things more appropriately.
	return fmt.Sprintf(`SOLANA_METRICS_CONFIG="%s"`, metricsConfig), nil
}

type InstallCommand struct {
	Flags    Flags
	KeyPairs KeyPairs
	Version  validator.Version
	Variant  *Variant
	Metrics  *Metrics
}

func (cmd *InstallCommand) Env() map[string]string {
	validatorEnv := ValidatorEnv{
		Metrics: cmd.Metrics,
	}

	env := map[string]string{
		"VALIDATOR_FLAGS":      strings.Join(cmd.Flags.ToArgs(), " "),
		"IDENTITY_KEYPAIR":     cmd.KeyPairs.Identity,
		"VOTE_ACCOUNT_KEYPAIR": cmd.KeyPairs.VoteAccount,
		"VALIDATOR_ENV":        validatorEnv.ToString(),
	}

	if cmd.Version != nil {
		env["VALIDATOR_VERSION"] = *cmd.Version
	}

	if cmd.Variant != nil {
		env["VALIDATOR_VARIANT"] = string(*cmd.Variant)
	} else {
		env["VALIDATOR_VARIANT"] = string(VariantAgave)
	}

	return env
}

func (cmd *InstallCommand) Script() string {
	return InstallScript
}

type ValidatorPaths struct {
	Accounts string `pulumi:"accounts"`
	Ledger   string `pulumi:"ledger"`
	Log      string `pulumi:"log"`
}

type Agave struct {
	Version  validator.Version `pulumi:"version,optional"`
	Variant  *Variant          `pulumi:"variant,optional"`
	KeyPairs KeyPairs          `pulumi:"keyPairs" provider:"secret"`
	Flags    Flags             `pulumi:"flags"`
	Metrics  *Metrics          `pulumi:"metrics,optional"`
}

func (agave *Agave) Install() runner.Command {
	return &InstallCommand{
		Flags:    agave.Flags,
		KeyPairs: agave.KeyPairs,
		Version:  agave.Version,
		Variant:  agave.Variant,
		Metrics:  agave.Metrics,
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
	var l []string

	// Note: These locations are hard coded inside asset-builder.
	l = append(l, f.S("identity", "/home/sol/validator-keypair.json"))
	l = append(l, f.S("vote-account", "/home/sol/vote-account-keypair.json"))

	if f.EntryPoint != nil {
		for _, entrypoint := range *f.EntryPoint {
			l = append(l, f.S("entrypoint", entrypoint))
		}
	}

	if f.KnownValidator != nil {
		for _, knownValidator := range *f.KnownValidator {
			l = append(l, f.S("known-validator", knownValidator))
		}
	}

	if f.ExpectedGenesisHash != nil {
		l = append(l, f.S("expected-genesis-hash", *f.ExpectedGenesisHash))
	}
	l = append(l, f.S("use-snapshot-archives-at-startup", f.UseSnapshotArchivesAtStartup))
	l = append(l, f.S("rpc-port", f.RpcPort))
	l = append(l, f.S("dynamic-port-range", f.DynamicPortRange))

	if f.GossipHost != nil {
		l = append(l, f.S("gossip-host", *f.GossipHost))
	}

	l = append(l, f.S("gossip-port", f.GossipPort))
	l = append(l, f.S("rpc-bind-address", f.RpcBindAddress))
	l = append(l, f.S("wal-recovery-mode", f.WalRecoveryMode))
	l = append(l, f.S("log", logPath))
	l = append(l, f.S("accounts", accountsPath))
	l = append(l, f.S("ledger", ledgerPath))
	l = append(l, f.S("limit-ledger-size", f.LimitLedgerSize))
	l = append(l, f.S("block-production-method", f.BlockProductionMethod))
	if f.TvuReceiveThreads != nil {
		l = append(l, f.S("tvu-receive-threads", *f.TvuReceiveThreads))
	}
	l = append(l, f.S("full-snapshot-interval-slots", f.FullSnapshotIntervalSlots))
	l = append(l, f.B("no-wait-for-vote-to-start-leader", f.NoWaitForVoteToStartLeader))
	l = append(l, f.B("only-known-rpc", f.OnlyKnownRPC))
	l = append(l, f.B("private-rpc", f.PrivateRPC))

	if f.FullRpcAPI != nil {
		l = append(l, f.B("full-rpc-api", *f.FullRpcAPI))
	}

	if f.NoVoting != nil {
		l = append(l, f.B("no-voting", *f.NoVoting))
	}

	if f.AllowPrivateAddr != nil {
		l = append(l, f.B("allow-private-addr", *f.AllowPrivateAddr))
	}

	if f.ExtraFlags != nil {
		l = append(l, *f.ExtraFlags...)
	}

	return l
}

func (Flags) S(k string, v interface{}) string {
	return fmt.Sprintf("--%s %v", k, v)
}

func (Flags) B(k string, v bool) string {
	if v {
		return fmt.Sprintf("--%s", k)
	}
	return ""
}

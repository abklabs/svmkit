package agave

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/solana"
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
	VariantXen         Variant = "xen"
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
		{
			Name:        string(VariantXen),
			Value:       VariantXen,
			Description: "The Xen validator",
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
		"VALIDATOR_FLAGS": strings.Join(cmd.Flags.Args(), " "),
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

		b.SetArray("SOLANA_CLI_CONFIG_FLAGS", conf.Flags().Args())
	}

	if t := cmd.TimeoutConfig; t != nil {
		b.Merge(t.Env())
	}

	b.SetP("VALIDATOR_VERSION", cmd.Version)

	if cmd.Variant != nil {
		b.Set("VALIDATOR_VARIANT", string(*cmd.Variant))
	} else {
		b.Set("VALIDATOR_VARIANT", string(VariantAgave))
	}

	b.Set("RPC_BIND_ADDRESS", cmd.Flags.RpcBindAddress)
	b.SetInt("RPC_PORT", cmd.Flags.RpcPort)

	if cmd.Flags.FullRpcAPI != nil && *cmd.Flags.FullRpcAPI && cmd.StartupPolicy != nil {
		b.SetBoolP("WAIT_FOR_RPC_HEALTH", cmd.StartupPolicy.WaitForRPCHealth)
	}

	if i := cmd.Info; i != nil {
		b.Set("VALIDATOR_INFO_NAME", i.Name)
		b.SetP("VALIDATOR_INFO_WEBSITE", i.Website)
		b.SetP("VALIDATOR_INFO_ICON_URL", i.IconURL)
		b.SetP("VALIDATOR_INFO_DETAILS", i.Details)
	}

	if s := cmd.ShutdownPolicy; s != nil {
		b.SetArray("VALIDATOR_EXIT_FLAGS", s.Flags().Args())
	}

	b.Set("LEDGER_PATH", ledgerPath)

	return b
}

func (cmd *InstallCommand) AddToPayload(p *runner.Payload) error {
	err := p.AddTemplate("steps.sh", installScriptTmpl, cmd)

	if err != nil {
		return err
	}

	if err := p.AddTemplate("check-validator", checkValidatorScriptTmpl, cmd); err != nil {
		return err
	}

	p.AddString("validator-keypair.json", cmd.KeyPairs.Identity)
	p.AddString("vote-account-keypair.json", cmd.KeyPairs.VoteAccount)

	return nil
}

type Agave struct {
	Environment    *solana.Environment   `pulumi:"environment,optional"`
	Version        *string               `pulumi:"version,optional"`
	Variant        *Variant              `pulumi:"variant,optional"`
	KeyPairs       KeyPairs              `pulumi:"keyPairs"`
	Flags          Flags                 `pulumi:"flags"`
	Metrics        *Metrics              `pulumi:"metrics,optional"`
	Info           *solana.ValidatorInfo `pulumi:"info,optional"`
	TimeoutConfig  *TimeoutConfig        `pulumi:"timeoutConfig,optional"`
	StartupPolicy  *StartupPolicy        `pulumi:"startupPolicy,optional"`
	ShutdownPolicy *ShutdownPolicy       `pulumi:"shutdownPolicy,optional"`
}

func (agave *Agave) Install() runner.Command {
	return &InstallCommand{
		Agave: *agave,
	}
}

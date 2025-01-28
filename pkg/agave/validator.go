package agave

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/solana"
)

const (
	accountsPath = "/home/sol/accounts"
	ledgerPath   = "/home/sol/ledger"
	logPath      = "/home/sol/log"

	identityKeyPairPath    = "/home/sol/validator-keypair.json"
	voteAccountKeyPairPath = "/home/sol/vote-account-keypair.json"
)

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

	packageInfo *PackageInfo
}

func (cmd *InstallCommand) Check() error {
	if m := cmd.Metrics; m != nil {
		if err := m.Check(); err != nil {
			return fmt.Errorf("Warning: Invalid metrics URL: %v\n", err)
		}
	}

	cmd.RunnerCommand.SetConfigDefaults()

	packageInfo, err := GeneratePackageInfo(cmd.Variant, cmd.Version)

	if err != nil {
		return err
	}

	if err := cmd.RunnerCommand.UpdatePackageGroup(packageInfo.PackageGroup); err != nil {
		return err
	}

	cmd.packageInfo = packageInfo

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

	b.Set("VALIDATOR_VARIANT", string(cmd.packageInfo.Variant))
	b.Set("VALIDATOR_PROCESS", cmd.packageInfo.Variant.ProcessName())
	b.Set("VALIDATOR_PACKAGE", cmd.packageInfo.Variant.PackageName())
	b.Merge(cmd.RunnerCommand.Env())

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

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	p.AddString("validator-keypair.json", cmd.KeyPairs.Identity)
	p.AddString("vote-account-keypair.json", cmd.KeyPairs.VoteAccount)

	return nil
}

type Agave struct {
	runner.RunnerCommand

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

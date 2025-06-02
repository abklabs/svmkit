package agave

import (
	_ "embed"
	"fmt"
	"net"
	"strings"

	"github.com/abklabs/svmkit/pkg/agave/geyser"
	"github.com/abklabs/svmkit/pkg/deletion"
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"github.com/abklabs/svmkit/pkg/solana"
	"github.com/abklabs/svmkit/pkg/validator"
)

const (
	accountsPath = "/home/sol/accounts"
	ledgerPath   = "/home/sol/ledger"

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
			return fmt.Errorf("warning: invalid metrics URL: %v", err)
		}
	}

	if g := cmd.GeyserPlugin; g != nil {
		if err := g.Check(); err != nil {
			return fmt.Errorf("warning: invalid geyser plugin config: %v", err)
		}

	}

	cmd.SetConfigDefaults()

	packageInfo, err := GeneratePackageInfo(cmd.GetVariant(), cmd.Version)

	if err != nil {
		return err
	}

	cmd.Variant = &packageInfo.Variant

	if g := cmd.GeyserPlugin; g != nil {
		if g.YellowstoneGRPC != nil {
			yellowstoneVersion := g.YellowstoneGRPC.Version
			packageInfo.PackageGroup.Add(deb.Package{Name: "svmkit-yellowstone_grpc", Version: &yellowstoneVersion})
		}
	}

	if err := cmd.UpdatePackageGroup(packageInfo.PackageGroup); err != nil {
		return err
	}

	cmd.packageInfo = packageInfo

	policy := cmd.GetDeletionPolicy()
	if err := policy.Check(); err != nil {
		return err
	}

	cmd.DeletionPolicy = &policy

	return nil
}

func (cmd *InstallCommand) Env() *runner.EnvBuilder {
	validatorEnv := runner.NewEnvBuilder()

	if m := cmd.Metrics; m != nil {
		validatorEnv.Set("SOLANA_METRICS_CONFIG", m.String())
	}

	b := runner.NewEnvBuilder()

	// Add plugin flag that points to config file that get's created by the runner
	if g := cmd.GeyserPlugin; g != nil {
		if cmd.Flags.GeyserPluginConfig == nil {
			pluginConfig := []string{"/home/sol/geyser-config.json"}
			cmd.Flags.GeyserPluginConfig = &pluginConfig
		} else {
			pluginConfig := append(*cmd.Flags.GeyserPluginConfig, "/home/sol/geyser-config.json")
			cmd.Flags.GeyserPluginConfig = &pluginConfig
		}
	}

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
	b.Set("VALIDATOR_SERVICE", cmd.packageInfo.Variant.ServiceName())
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

	if g := cmd.GeyserPlugin; g != nil {
		if g.YellowstoneGRPC != nil {
			b.SetBool("YELLOWSTONE_GRPC", true)

			if g.YellowstoneGRPC.Config != nil {
				address := g.YellowstoneGRPC.Config.Grpc.Address
				_, port, err := net.SplitHostPort(address)
				if err == nil {
					b.Set("YELLOWSTONE_GRPC_PORT", port)
				}
			}
		}
	}

	cmd.DeletionPolicy.Create(&cmd.Agave, b)

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

	if plugin := cmd.GeyserPlugin; plugin != nil {
		confString, err := plugin.ToConfigString()
		if err != nil {
			return err
		}
		p.AddString("geyser-config.json", confString)
	}

	if err := deletion.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

type UninstallCommand struct {
	Agave
}

// AddToPayload implements runner.Command.
// Subtle: this method shadows the method (Agave).AddToPayload of UninstallCommand.Agave.
func (u *UninstallCommand) AddToPayload(p *runner.Payload) error {
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
func (u *UninstallCommand) Check() error {
	u.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup()

	if err := u.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	policy := u.GetDeletionPolicy()
	if err := policy.Check(); err != nil {
		return err
	}

	u.DeletionPolicy = &policy

	return nil
}

// Config implements runner.Command.
// Subtle: this method shadows the method (Agave).Config of UninstallCommand.Agave.
func (u *UninstallCommand) Config() *runner.Config {
	return u.RunnerConfig
}

// Env implements runner.Command.
// Subtle: this method shadows the method (Agave).Env of UninstallCommand.Agave.
func (u *UninstallCommand) Env() *runner.EnvBuilder {
	b := runner.NewEnvBuilder()

	b.Merge(u.RunnerCommand.Env())

	u.DeletionPolicy.Delete(&u.Agave, b)

	return b
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
	GeyserPlugin   *geyser.GeyserPlugin  `pulumi:"geyserPlugin,optional"`
	DeletionPolicy *deletion.Policy      `pulumi:"deletionPolicy,optional"`
}

func (agave *Agave) Install() runner.Command {
	return &InstallCommand{
		Agave: *agave,
	}
}

func (agave *Agave) GetVariant() Variant {
	if agave.Variant == nil {
		return VariantAgave
	} else {
		return *agave.Variant
	}
}

func (agave *Agave) GetDeletionPolicy() deletion.Policy {
	if agave.DeletionPolicy == nil {
		return deletion.PolicyKeep
	} else {
		return *agave.DeletionPolicy
	}
}

func (agave *Agave) Properties() validator.Properties {
	variant := agave.GetVariant()

	return validator.Properties{SystemdServiceName: variant.ServiceName()}
}

func (agave *Agave) Uninstall() runner.Command {
	return &UninstallCommand{
		Agave: *agave,
	}
}

func (agave *Agave) ManagedFiles() []string {
	return []string{accountsPath, ledgerPath}
}

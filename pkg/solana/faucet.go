package solana

import (
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
)

const (
	faucetKeyPairPath = "/home/sol/faucet-keypair.json"
)

type InstallCommand struct {
	Faucet
}

func (cmd *InstallCommand) Env() *runner.EnvBuilder {

	faucetEnv := runner.NewEnvBuilder()

	b := runner.NewEnvBuilder()

	b.SetMap(map[string]string{
		"FAUCET_FLAGS": strings.Join(cmd.Flags.ToArgs(), " "),
		"FAUCET_ENV":   faucetEnv.String(),
	})

	b.Set("FAUCET_PORT", "9900") // hardcoded in solana-faucet

	b.SetP("FAUCET_VERSION", cmd.Version)

	return b

}

func (cmd *InstallCommand) Check() error {
	return nil
}

func (cmd *InstallCommand) AddToPayload(p *runner.Payload) error {
	faucetScript, err := assets.Open(assetsFaucetScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", faucetScript)
	p.AddString("faucet-keypair.json", cmd.KeyPair)

	return nil
}

type Faucet struct {
	Flags   FaucetFlags `pulumi:"flags"`
	Version *string     `pulumi:"version,optional"`
	KeyPair string      `pulumi:"keypair" provider:"secret"`
}

func (f *Faucet) ToArgs() []string {
	return f.Flags.ToArgs()
}

func (f *Faucet) Install() runner.Command {
	return &InstallCommand{
		Faucet: *f,
	}
}

type FaucetFlags struct {
	// Optional. Allow requests from specified IPs without request limit.
	// If multiple --allow-ip flags are provided, all specified IPs are allowed.
	AllowIPs *[]string `pulumi:"allowIPs,optional"`

	// Optional. Request limit for a single request, in SOL.
	// If not specified, no limit is applied.
	PerRequestCap *int `pulumi:"perRequestCap,optional"`

	// Optional. Request limit for a given time slice, in SOL.
	// If not specified, no limit is applied.
	PerTimeCap *int `pulumi:"perTimeCap,optional"`

	// Optional. Length of the time slice in seconds.
	// If not specified, no slicing is applied.
	SliceSeconds *int `pulumi:"sliceSeconds,optional"`
}

func (f *FaucetFlags) ToArgs() []string {
	b := runner.FlagBuilder{}

	b.Append("--keypair", faucetKeyPairPath)

	if f.AllowIPs != nil {
		for _, ip := range *f.AllowIPs {
			b.AppendP("allow-ip", &ip)
		}
	}

	b.AppendIntP("per-request-cap", f.PerRequestCap)
	b.AppendIntP("per-time-cap", f.PerTimeCap)
	b.AppendIntP("slice-seconds", f.SliceSeconds)

	return b.ToArgs()
}

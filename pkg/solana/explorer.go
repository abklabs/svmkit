package solana

import (
	"strings"

	"github.com/abklabs/svmkit/pkg/runner"
)

const (
	defaultHostname = "0.0.0.0"
	defaultPort     = 3000
)

type ExplorerCommand struct {
	Explorer
}

func (cmd *ExplorerCommand) Env() *runner.EnvBuilder {

	explorerEnv := runner.NewEnvBuilder()

	if cmd.RPCURL != nil {
		explorerEnv.SetP("NEXT_PUBLIC_CLUSTER_URI", cmd.RPCURL)
	} else {
		explorerEnv.SetP("NEXT_PUBLIC_CLUSTER_URI", cmd.Environment.RPCURL)
	}

	explorerEnv.SetP("NEXT_PUBLIC_CLUSTER_NAME", cmd.ClusterName)
	explorerEnv.SetP("NEXT_PUBLIC_EXPLORER_NAME", cmd.Name)
	explorerEnv.SetP("NEXT_PUBLIC_EXPLORER_SYMBOL", cmd.Symbol)

	b := runner.NewEnvBuilder()

	b.SetMap(map[string]string{
		"EXPLORER_FLAGS": strings.Join(cmd.Flags.Args(), " "),
		"EXPLORER_ENV":   explorerEnv.String(),
	})

	b.SetIntP("EXPLORER_PORT", cmd.Flags.Port)
	b.SetP("EXPLORER_VERSION", cmd.Version)

	return b

}

func (cmd *ExplorerCommand) Check() error {
	return nil
}

func (cmd *ExplorerCommand) AddToPayload(p *runner.Payload) error {
	explorerScript, err := assets.Open(assetsExplorerScript)

	if err != nil {
		return err
	}

	p.AddReader("steps.sh", explorerScript)

	return nil
}

type Explorer struct {
	Environment Environment   `pulumi:"environment"`
	Flags       ExplorerFlags `pulumi:"flags"`
	Version     *string       `pulumi:"version,optional"`
	Name        *string       `pulumi:"name,optional"`
	Symbol      *string       `pulumi:"symbol,optional"`
	ClusterName *string       `pulumi:"clusterName,optional"`
	RPCURL      *string       `pulumi:"RPCURL,optional"`
}

func (f *Explorer) Install() runner.Command {
	return &ExplorerCommand{
		Explorer: *f,
	}
}

type ExplorerFlags struct {
	Hostname         *string `pulumi:"hostname,optional"`
	Port             *int    `pulumi:"port,optional"`
	KeepAliveTimeout *int    `pulumi:"keepAliveTimeout,optional"`
}

func (f *ExplorerFlags) Args() []string {
	b := runner.FlagBuilder{}

	if f.Hostname != nil {
		b.AppendP("hostname", f.Hostname)
	} else {
		value := defaultHostname
		b.AppendP("hostname", &value)
	}

	if f.Port != nil {
		b.AppendIntP("port", f.Port)
	} else {
		value := defaultPort
		b.AppendIntP("port", &value)
	}

	b.AppendIntP("keepAliveTimeout", f.KeepAliveTimeout)

	return b.Args()
}

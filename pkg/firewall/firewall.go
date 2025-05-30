package firewall

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"

	"dario.cat/mergo"
)

type FirewallCommand struct {
	Firewall
}

func (cmd *FirewallCommand) Check() error {
	cmd.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup("ufw")

	if err := cmd.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	return nil
}

func (cmd *FirewallCommand) AddToPayload(p *runner.Payload) error {
	if err := p.AddTemplate("steps.sh", firewallScriptTmpl, cmd); err != nil {
		return err
	}
	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

type FirewallParams struct {
	AllowPorts []string `pulumi:"allowPorts" toml:"allowPorts"`
}

type Firewall struct {
	runner.RunnerCommand
	Params FirewallParams `pulumi:"params" toml:"params"`
}

func (f *Firewall) Create() runner.Command {
	return &FirewallCommand{
		Firewall: *f,
	}
}

func (t *Firewall) Merge(other *Firewall) error {
	if other == nil {
		return nil
	}
	return mergo.Merge(t, other, mergo.WithOverride)
}

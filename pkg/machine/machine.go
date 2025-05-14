package machine

import (
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
)

type Machine struct {
	runner.RunnerCommand
}

type CreateCommand struct {
	Machine
}

func (cmd *CreateCommand) Env() *runner.EnvBuilder {
	tunerEnv := runner.NewEnvBuilder()
	tunerEnv.Merge(cmd.RunnerCommand.Env())
	return tunerEnv
}

func (cmd *CreateCommand) Check() error {
	cmd.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup()

	if err := cmd.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	return nil
}

func (cmd *CreateCommand) AddToPayload(p *runner.Payload) error {
	if err := p.AddTemplate("steps.sh", installScriptTmpl, cmd); err != nil {
		return err
	}

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

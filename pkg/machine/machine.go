package machine

import (
	"github.com/abklabs/svmkit/pkg/machine/apt"
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
)

type Machine struct {
	runner.RunnerCommand

	AptConfig *apt.Config `pulumi:"aptConfig,optional"`
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
	if err := p.AddTemplate(runner.ScriptNameSteps, installScriptTmpl, cmd); err != nil {
		return err
	}

	sources := apt.Sources{}

	excludeDefaultSources := false

	if conf := cmd.AptConfig; conf != nil {
		if conf.ExcludeDefaultSources != nil {
			excludeDefaultSources = *conf.ExcludeDefaultSources
		}

		if conf.Sources != nil {
			sources = append(sources, *conf.Sources...)
		}
	}

	if !excludeDefaultSources {
		// Attach our default apt repo
		sources = append(sources, apt.Source{
			Types:      []string{"deb"},
			URIs:       []string{"https://apt.abklabs.com/svmkit"},
			Suites:     []string{"dev "},
			Components: []string{"main"},
			SignedBy: &apt.SignedBy{
				PublicKey: &ABKLabsArchiveDevPubKey,
			},
		})
	}

	res, err := sources.MarshalText()

	if err != nil {
		return err
	}

	p.NewBuffer(runner.PayloadFile{Path: "svmkit.sources"}, res)

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

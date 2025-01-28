package runner

import (
	"github.com/abklabs/svmkit/pkg/runner/deb"
)

type RunnerCommand struct {
	RunnerConfig *Config `pulumi:"runnerConfig,optional"`

	packageGroup *deb.PackageGroup
}

func (r *RunnerCommand) UpdatePackageGroup(grp *deb.PackageGroup) error {
	if r.RunnerConfig != nil {
		if err := r.RunnerConfig.UpdatePackageGroup(grp); err != nil {
			return err
		}
	}

	r.packageGroup = grp

	return nil
}

func (r *RunnerCommand) AddToPayload(p *Payload) error {
	if r.packageGroup == nil {
		panic("payload cannot be added if the package group hasn't been updated!")
	}

	return r.packageGroup.AddToPayload(p)
}

func (r *RunnerCommand) Env() *EnvBuilder {
	if r.packageGroup == nil {
		panic("environment cannot be configured if the package group hasn't been updated!")
	}

	env := NewEnvBuilder()
	env.SetArray("PACKAGE_LIST", r.packageGroup.Args())
	return env
}

package firedancer

import (
	"github.com/abklabs/svmkit/pkg/runner"
)

type KeyPairs struct {
	Identity    string `pulumi:"identity" provider:"secret"`
	VoteAccount string `pulumi:"voteAccount" provider:"secret"`
}

type Firedancer struct {
	KeyPairs KeyPairs `pulumi:"keyPairs"`
	Config   Config   `pulumi:"config"`
}

func (fd *Firedancer) Install() runner.Command {
	return &InstallCommand{
		Firedancer: *fd,
	}
}

type InstallCommand struct {
	Firedancer
}

func (c *InstallCommand) Check() error {
	return nil
}

func (c *InstallCommand) Env() *runner.EnvBuilder {
	e := runner.NewEnvBuilder()
	return e
}

func (c *InstallCommand) AddToPayload(p *runner.Payload) error {
	{
		w := p.NewWriter(runner.PayloadFile{Path: "config.toml"})

		if err := c.Config.Encode(w); err != nil {
			return err
		}
	}

	{
		r, err := assets.Open(assetsInstall)

		if err != nil {
			return err
		}

		p.AddReader("steps.sh", r)
	}

	{
		r, err := assets.Open(assetsFDService)

		if err != nil {
			return err
		}

		p.AddReader("svmkit-fd-validator.service", r)
	}

	{
		r, err := assets.Open(assetsFDSetupService)

		if err != nil {
			return err
		}

		p.AddReader("svmkit-fd-setup.service", r)
	}

	p.AddString("validator-keypair.json", c.KeyPairs.Identity)
	p.AddString("vote-account-keypair.json", c.KeyPairs.VoteAccount)

	return nil
}

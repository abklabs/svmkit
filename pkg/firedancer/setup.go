package firedancer

import (
	"fmt"

	"github.com/abklabs/svmkit/pkg/deletion"
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"github.com/abklabs/svmkit/pkg/solana"
	"github.com/abklabs/svmkit/pkg/validator"
)

const (
	identityKeyPairPath = "/home/sol/validator-keypair.json"
)

type KeyPairs struct {
	Identity    string `pulumi:"identity" provider:"secret"`
	VoteAccount string `pulumi:"voteAccount" provider:"secret"`
}

type Firedancer struct {
	runner.RunnerCommand

	Environment    *solana.Environment `pulumi:"environment,optional"`
	Version        *string             `pulumi:"version,optional"`
	Variant        *Variant            `pulumi:"variant,optional"`
	DeletionPolicy *deletion.Policy    `pulumi:"deletionPolicy,optional"`

	KeyPairs KeyPairs `pulumi:"keyPairs"`
	Config   Config   `pulumi:"config"`
}

func (fd *Firedancer) Install() runner.Command {
	return &InstallCommand{
		Firedancer: *fd,
	}
}

func (fd *Firedancer) Uninstall() runner.Command {
	return &UninstallCommand{
		Firedancer: *fd,
	}
}

func (fd *Firedancer) GetVariant() Variant {
	if fd.Variant == nil {
		return VariantFrankendancer
	} else {
		return *fd.Variant
	}
}

func (fd *Firedancer) GetDeletionPolicy() deletion.Policy {
	if fd.DeletionPolicy == nil {
		return deletion.PolicyKeep
	} else {
		return *fd.DeletionPolicy
	}
}

func (fd *Firedancer) Properties() validator.Properties {
	variant := fd.GetVariant()
	return validator.Properties{SystemdServiceName: variant.ServiceName()}
}

func (fd *Firedancer) checkManagedFiles() error {
	if fd.Config.Ledger.Path == nil {
		return fmt.Errorf("missing ledger.path in Firedancer config")
	}

	if fd.Config.Ledger.AccountsPath == nil {
		return fmt.Errorf("missing ledger.accountsPath in Firedancer config")
	}

	return nil
}

func (fd *Firedancer) ManagedFiles() []string {
	return []string{*fd.Config.Ledger.Path, *fd.Config.Ledger.AccountsPath}
}

type InstallCommand struct {
	Firedancer
}

func (c *InstallCommand) Check() error {
	c.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup("svmkit-solana-cli")
	pkgGrp.Add(deb.Package{Name: "svmkit-frankendancer", Version: c.Version})

	variant := c.GetVariant()

	if err := variant.Check(); err != nil {
		return err
	}

	c.Variant = &variant

	if err := c.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	if c.DeletionPolicy != nil {
		if err := c.checkManagedFiles(); err != nil {
			return err
		}
	}

	policy := c.GetDeletionPolicy()
	if err := policy.Check(); err != nil {
		return err
	}

	c.DeletionPolicy = &policy

	return nil
}

func (c *InstallCommand) Env() *runner.EnvBuilder {
	e := runner.NewEnvBuilder()

	{
		s := identityKeyPairPath
		conf := solana.CLIConfig{
			KeyPair: &s,
		}

		if senv := c.Environment; senv != nil {
			conf.URL = senv.RPCURL
		}

		e.SetArray("SOLANA_CLI_CONFIG_FLAGS", conf.Flags().Args())
	}

	e.Merge(c.RunnerCommand.Env())
	e.Set("VALIDATOR_PACKAGE", c.Variant.PackageName())
	e.Set("VALIDATOR_SERVICE", c.Variant.ServiceName())

	c.DeletionPolicy.Create(&c.Firedancer, e)

	return e
}

func (c *InstallCommand) AddToPayload(p *runner.Payload) error {
	{
		w := p.NewWriter(runner.PayloadFile{Path: "config.toml"})

		if err := c.Firedancer.Config.Encode(w); err != nil {
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

		p.AddReader(fmt.Sprintf("%s.service", c.Variant.ServiceName()), r)
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

	if err := c.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	if err := deletion.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

func (c *InstallCommand) Config() *runner.Config {
	return c.RunnerConfig
}

type UninstallCommand struct {
	Firedancer
}

func (u *UninstallCommand) Check() error {
	u.SetConfigDefaults()

	variant := u.GetVariant()

	if err := variant.Check(); err != nil {
		return err
	}

	u.Variant = &variant

	pkgGrp := deb.Package{}.MakePackageGroup("svmkit-solana-cli")
	pkgGrp.Add(deb.Package{Name: "svmkit-frankendancer", Version: u.Version})

	if err := u.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	if u.DeletionPolicy != nil {
		if err := u.checkManagedFiles(); err != nil {
			return err
		}
	}

	policy := u.GetDeletionPolicy()
	if err := policy.Check(); err != nil {
		return err
	}

	u.DeletionPolicy = &policy

	return nil
}

func (u *UninstallCommand) Env() *runner.EnvBuilder {
	e := runner.NewEnvBuilder()

	e.Merge(u.RunnerCommand.Env())
	e.Set("VALIDATOR_PACKAGE", u.Variant.PackageName())
	e.Set("VALIDATOR_SERVICE", u.Variant.ServiceName())

	u.DeletionPolicy.Delete(&u.Firedancer, e)

	return e
}

func (u *UninstallCommand) AddToPayload(p *runner.Payload) error {
	{
		r, err := assets.Open(assetsUninstall)

		if err != nil {
			return err
		}

		p.AddReader("steps.sh", r)
	}

	if err := deletion.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

func (u *UninstallCommand) Config() *runner.Config {
	return u.RunnerConfig
}

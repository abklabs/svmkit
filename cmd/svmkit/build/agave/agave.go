package agave

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"github.com/abklabs/svmkit/pkg/runner/deployer"

	"github.com/abklabs/svmkit/cmd/svmkit/build/agave/assets"
	"github.com/abklabs/svmkit/cmd/svmkit/utils"
)

type Build struct {
	runner.RunnerCommand

	BuildDir           string
	Maintainer         string
	UseAlterativeClang bool
	BuildExtras        bool
	NoBuild            bool

	ValidatorTarget string
	PackagePrefix   string
}

func (cmd *Build) Env() *runner.EnvBuilder {
	env := runner.NewEnvBuilder()

	env.Set("MAINTAINER", cmd.Maintainer)
	env.Set("PACKAGE_PREFIX", cmd.PackagePrefix)
	env.Set("TARGET_VALIDATOR", cmd.ValidatorTarget)
	env.Set("BUILD_DIR", cmd.BuildDir)
	env.SetBool("USE_ALTERNATIVE_CLANG", cmd.UseAlterativeClang)
	env.SetBool("BUILD_EXTRAS", cmd.BuildExtras)
	env.SetBool("NO_BUILD", cmd.NoBuild)

	return env
}

func (cmd *Build) Check() error {
	cmd.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup()

	if err := cmd.UpdatePackageGroup(pkgGrp); err != nil {
		return err
	}

	return nil
}

func (cmd *Build) AddToPayload(p *runner.Payload) error {
	script, err := assets.FS.Open(assets.BuildScript)

	if err != nil {
		return err
	}

	reader, err := utils.AssembleScript(AgaveCmd.Flags(), script)

	if err != nil {
		return err
	}

	p.AddReader(runner.ScriptNameSteps, reader)

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

var AgaveCmd = &cobra.Command{
	Use:   "agave",
	Short: "Build an Agave Debian package",
	Long:  "Build an Agave (or Agave-variant) Debian package",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()

		if err != nil {
			return err
		}

		flags := cmd.Flags()

		maintainer, err := flags.GetString("maintainer")

		if err != nil {
			return err
		}

		useAlternativeClang, err := flags.GetBool("use-alternative-clang")

		if err != nil {
			return err
		}

		buildExtras, err := flags.GetBool("build-extras")

		if err != nil {
			return err
		}

		noBuild, err := flags.GetBool("no-build")

		if err != nil {
			return err
		}

		validatorTarget, err := flags.GetString("validator-target")

		if err != nil {
			return err
		}

		packagePrefix, err := flags.GetString("package-prefix")

		if err != nil {
			return err
		}

		keepPayload, err := flags.GetBool("keep-payload")

		if err != nil {
			return err
		}

		dryRun, err := flags.GetBool("dry-run")

		if err != nil {
			return err
		}

		runnerCommand := &Build{
			BuildDir:           cwd,
			Maintainer:         maintainer,
			UseAlterativeClang: useAlternativeClang,
			BuildExtras:        buildExtras,
			NoBuild:            noBuild,
			ValidatorTarget:    validatorTarget,
			PackagePrefix:      packagePrefix,
		}

		if err := runnerCommand.Check(); err != nil {
			return err
		}

		outputDir, err := os.MkdirTemp("", "build-*")

		if err != nil {
			return err
		}

		p := &runner.Payload{
			RootPath:    outputDir,
			DefaultMode: 0640,
		}

		if err := runner.PrepareCommandPayload(p, runnerCommand); err != nil {
			return err
		}

		d := &deployer.Local{
			Payload:     p,
			KeepPayload: keepPayload,
		}

		log.Printf("writing to '%s'...", outputDir)

		if err := d.Deploy(); err != nil {
			return err
		}

		handler := &deployer.LoggerHandler{
			LogCallback: func(s string) {
				log.Print(s)
			},
		}

		if dryRun {
			log.Printf("performing a dry run; no commands executed")
			return nil
		}

		return d.Run([]string{"./run.sh"}, handler)
	},
}

func init() {
	flags := AgaveCmd.Flags()

	flags.String("maintainer", "Engineering <engineering@abklabs.com>", "name and email of the maintainer of the package")
	flags.String("validator-target", "agave-validator", "cargo build target to use for the validator")
	flags.String("package-prefix", "svmkit-agave", "prefix to use with built packages")
	flags.Bool("build-extras", false, "should build extra packages (e.g. Solana CLI)")
	flags.Bool("use-alternative-clang", false, "use an older clang (e.g. 14) for the build")
	flags.Bool("keep-payload", false, "don't remove build scripts after completion")
	flags.Bool("no-build", false, "configure the repository, but do not build")
}

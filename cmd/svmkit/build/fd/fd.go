package fd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"github.com/abklabs/svmkit/pkg/runner/deployer"

	"github.com/abklabs/svmkit/cmd/svmkit/build/fd/assets"
	"github.com/abklabs/svmkit/cmd/svmkit/utils"
)

type Build struct {
	runner.RunnerCommand

	BuildDir       string
	KeepPayload    bool
	DepsFetchExtra string
	MakeMachine    string
	MakeCFlags     string
	MakeTarget     string
}

func (cmd *Build) Env() *runner.EnvBuilder {
	env := runner.NewEnvBuilder()

	env.Set("BUILD_DIR", cmd.BuildDir)
	env.Set("FD_DEPS_FETCH_EXTRA", cmd.DepsFetchExtra)
	env.Set("FD_MAKE_MACHINE", cmd.MakeMachine)
	env.Set("FD_MAKE_CFLAGS", cmd.MakeCFlags)
	env.Set("FD_MAKE_TARGET", cmd.MakeTarget)

	return env
}

func (cmd *Build) Check() error {
	cmd.SetConfigDefaults()

	pkgGrp := deb.Package{}.MakePackageGroup("cmake", "alien")

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

	reader, err := utils.AssembleScript(FDCmd.Flags(), script)

	if err != nil {
		return err
	}

	p.AddReader(runner.ScriptNameSteps, reader)

	if err := cmd.RunnerCommand.AddToPayload(p); err != nil {
		return err
	}

	return nil
}

var FDCmd = &cobra.Command{
	Use:   "fd",
	Short: "Build a Frankendancer Debian package",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		flags := cmd.Flags()

		keepPayload, err := flags.GetBool("keep-payload")
		if err != nil {
			return err
		}

		depsFetchExtra, err := flags.GetString("deps-fetch-extra")
		if err != nil {
			return err
		}

		makeMachine, err := flags.GetString("make-machine")
		if err != nil {
			return err
		}

		makeCFlags, err := flags.GetString("make-cflags")
		if err != nil {
			return err
		}

		makeTarget, err := flags.GetString("make-target")
		if err != nil {
			return err
		}

		dryRun, err := flags.GetBool("dry-run")

		if err != nil {
			return err
		}

		runnerCommand := &Build{
			BuildDir:       cwd,
			KeepPayload:    keepPayload,
			MakeMachine:    makeMachine,
			MakeCFlags:     makeCFlags,
			DepsFetchExtra: depsFetchExtra,
			MakeTarget:     makeTarget,
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
	flags := FDCmd.Flags()

	flags.Bool("keep-payload", false, "don't remove build scripts after completion")
	flags.String("deps-fetch-extra", "+dev", "extra deps.sh fetch args")
	flags.String("make-machine", "linux_gcc_x86_64", "MACHINE env var passed to make")
	flags.String("make-cflags", "", "extra make CFLAGS")
	flags.String("make-target", "fdctl", "TARGET env var passed to make")
}

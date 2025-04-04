package fd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deb"
	"github.com/abklabs/svmkit/pkg/runner/deployer"

	"github.com/abklabs/svmkit/cmd/svmkit/build/fd/assets"
)

type Build struct {
	runner.RunnerCommand

	BuildDir    string
	KeepPayload bool
}

func (cmd *Build) Env() *runner.EnvBuilder {
	env := runner.NewEnvBuilder()

	env.Set("BUILD_DIR", cmd.BuildDir)

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

	p.AddReader("steps.sh", script)

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

		runnerCommand := &Build{
			BuildDir:    cwd,
			KeepPayload: keepPayload,
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

		return d.Run([]string{"./run.sh"}, handler)
	},
}

func init() {
	flags := FDCmd.Flags()

	flags.Bool("keep-payload", false, "don't remove build scripts after completion")
}

package component

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"

	"github.com/abklabs/svmkit/pkg/registry"
	"github.com/abklabs/svmkit/pkg/runner"
	"github.com/abklabs/svmkit/pkg/runner/deployer"
)

var GenerateCmd = &cobra.Command{
	Use: "generate",
}

func makeCommandGlue(runnerCommand runner.Command) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		inputFilename := args[0]
		outputDir := args[1]

		_, err := toml.DecodeFile(inputFilename, runnerCommand)

		if err != nil {
			return err
		}

		err = runnerCommand.Check()

		if err != nil {
			return err
		}

		p := &runner.Payload{
			RootPath:    outputDir,
			DefaultMode: 0640,
		}

		err = runner.PrepareCommandPayload(p, runnerCommand)

		if err != nil {
			return err
		}

		d := &deployer.Local{
			Payload:     p,
			KeepPayload: true,
		}

		log.Printf("using '%s' as inputs, writing to '%s'...", inputFilename, outputDir)

		return d.Deploy()
	}
}

func init() {
	for _, comp := range registry.Components {
		compCommand := &cobra.Command{
			Use:   comp.Name.String(),
			Short: comp.Summary,
		}

		for _, op := range comp.Op {
			cmd := &cobra.Command{
				Use:  fmt.Sprintf("%s inputTOML outputDir", op.Action.String()),
				RunE: makeCommandGlue(op.Creator()),
				Args: cobra.ExactArgs(2),
			}

			compCommand.AddCommand(cmd)
		}

		GenerateCmd.AddCommand(compCommand)
	}
}

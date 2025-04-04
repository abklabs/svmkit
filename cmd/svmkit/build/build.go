package build

import (
	"github.com/abklabs/svmkit/cmd/svmkit/build/agave"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use: "build",
}

func init() {
	BuildCmd.AddCommand(agave.AgaveCmd)
}

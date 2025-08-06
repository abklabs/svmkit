package build

import (
	"github.com/abklabs/svmkit/cmd/svmkit/build/agave"
	"github.com/abklabs/svmkit/cmd/svmkit/build/fd"
	"github.com/abklabs/svmkit/cmd/svmkit/build/yellowstone_grpc"
	"github.com/abklabs/svmkit/cmd/svmkit/utils"
	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Debian Packages of Solana tooling",
}

func init() {
	utils.AddScriptFlags(BuildCmd.PersistentFlags())
	BuildCmd.AddCommand(agave.AgaveCmd)
	BuildCmd.AddCommand(yellowstone_grpc.YellowstoneGRPCCmd)
	BuildCmd.AddCommand(fd.FDCmd)
}

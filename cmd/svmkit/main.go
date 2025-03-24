package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/abklabs/svmkit/cmd/svmkit/component"
)

var (
	rootCmd = &cobra.Command{
		Use:   "svmkit",
		Short: "A CLI for interacting with svmkit",
		Long:  `svmkit is a wrapper around the library of the same name.`,
	}
)

func init() {
	rootCmd.AddCommand(component.GenerateCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

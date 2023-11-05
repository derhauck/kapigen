package cmd

import (
	"github.com/spf13/cobra"
	"kapigen.kateops.com/cmd/generate"
	"kapigen.kateops.com/cmd/version"
	"os"
)

var Verbose bool
var rootCmd = &cobra.Command{
	Short:            "Kateops Pipeline Generator",
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(generate.Cmd)
	rootCmd.AddCommand(version.Cmd)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "will log verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

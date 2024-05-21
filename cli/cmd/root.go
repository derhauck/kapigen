package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Verbose bool
var rootCmd = &cobra.Command{
	Short:            "Kateops Pipeline Generator",
	TraverseChildren: true,
}

func init() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "will log verbose output")
	rootCmd.PersistentFlags().String("private-token", "", "ENV var name to use for private token")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

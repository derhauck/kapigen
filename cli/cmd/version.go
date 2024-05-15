package cmd

import (
	"github.com/spf13/cobra"
	cmd "kapigen.kateops.com/cmd/version"
)

var versionCmd = &cobra.Command{
	Use:              "version",
	Short:            "Will allow modification or creation of new version",
	Long:             "The version command focuses on the modification and creation of semantic versions for either file based versioning or via gitlab tags",
	TraverseChildren: true,
}

func init() {
	versionCmd.AddCommand(cmd.NewCmd)
}

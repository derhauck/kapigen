package cmd

import (
	"github.com/spf13/cobra"
	"kapigen.kateops.com/cmd/generate/pipeline"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate resource",
}

func init() {
	generateCmd.AddCommand(pipeline.Cmd)
}

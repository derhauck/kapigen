package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/kateops/kapigen/cli/cmd/pipeline"
)

var pipelineCmd = &cobra.Command{
	Use:   "pipeline",
	Short: "Pipeline resource",
}

func init() {
	pipelineCmd.AddCommand(pipeline.GenerateCmd)
	pipelineCmd.AddCommand(pipeline.ReportsCmd)
}

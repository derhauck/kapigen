package generate

import (
	"github.com/spf13/cobra"
	"kapigen.kateops.com/cmd/generate/pipeline"
)

var Cmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate resource",
}

func init() {
	Cmd.AddCommand(pipeline.Cmd)
}

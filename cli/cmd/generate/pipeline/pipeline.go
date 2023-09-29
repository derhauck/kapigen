package pipeline

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/config"
	"os"
)

var Cmd = &cobra.Command{
	Use:              "pipeline",
	Short:            "Generate pipeline",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.NewPersistentConfig(cmd)
		logger.Debug("activated verbose mode")
		configPath := "test.yaml"
		body, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		var pipelineConfig config.PipelineConfig
		err = yaml.Unmarshal(body, &pipelineConfig)
		if err != nil {
			return err
		}
		for i := 0; i < len(pipelineConfig.Pipelines); i++ {
			configuration := pipelineConfig.Pipelines[i]
			err = configuration.Decode()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	Cmd.Flags().Bool("file", false, "output file")
	Cmd.Flags().Bool("config", false, "config to use")

}

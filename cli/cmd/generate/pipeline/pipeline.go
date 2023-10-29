package pipeline

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/gitlab/pipeline"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/config"
	"kapigen.kateops.com/internal/pipeline/types"
	"os"
)

var Cmd = &cobra.Command{
	Use:              "pipeline",
	Short:            "Generate pipeline",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.NewPersistentConfig(cmd)
		logger.Debug("activated verbose mode")

		configPath, err := cmd.Flags().GetString("config")
		if err != nil {
			return err
		}
		pipelineFile, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		logger.Info("will read pipeline config from: " + configPath)
		pipelineJobs, err := types.LoadJobsFromPipelineConfig(factory.New(), configPath, config.PipelineConfigTypes)
		if err != nil {
			return err
		}
		logger.Info("ci jobs created")
		merge, err := cmd.Flags().GetBool("no-merge")
		if err != nil {
			return err
		}
		var ciPipeline map[string]interface{}
		if merge == false {
			pipelineJobs, err = pipelineJobs.DynamicMerge()
			if err != nil {
				return err
			}
			logger.Info("ci jobs dynamically merged")
		}

		pipelineJobs, err = pipelineJobs.EvaluateNames()
		if err != nil {
			return err
		}
		logger.Info("ci jobs named to be unique")
		ciPipeline = types.JobsToMap(pipelineJobs)
		logger.Info("ci job list converted to map")
		pipeline.NewDefaultCiPipeline().Render().AddToMap(ciPipeline)
		logger.Info("ci jobs rendered")

		data, err := yaml.Marshal(ciPipeline)
		if err != nil {
			return err
		}
		logger.Info("converted pipeline to yaml")

		err = os.WriteFile(pipelineFile, data, 0777)
		logger.Info("wrote yaml to file: " + pipelineFile)

		return err
	},
}

func init() {
	Cmd.Flags().String("file", "pipeline.yaml", "output file")
	Cmd.Flags().String("config", "config.kapigen.yaml", "config to use")
	Cmd.Flags().Bool("no-merge", false, "use dynamic job merge")

}

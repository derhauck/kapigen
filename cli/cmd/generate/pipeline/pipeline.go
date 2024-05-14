package pipeline

import (
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/gitlab/pipeline"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/config"
	"kapigen.kateops.com/internal/pipeline/jobs"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/version"
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
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			return err
		}
		logger.Info("will create settings")
		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(mode)),
		)

		logger.Info("will read pipeline config from: " + configPath)
		pipelineJobs, pipelineConfig, err := types.LoadJobsFromPipelineConfig(factory.New(settings), configPath, config.PipelineConfigTypes)
		if err != nil {
			return err
		}
		logger.Info("ci jobs created")

		if pipelineConfig.Noop {
			logger.Info("noop mode activated, will add \"Noop\" job to pipeline")
			pipelineJobs.AddJob(jobs.NewNoop())
		}

		if pipelineConfig.Tag {
			logger.Info("tag mode activated, will add \"Tag\" job to pipeline")
			pipelineJobs.AddJob(jobs.NewTag()).
				AddJob(jobs.NewTagKapigen())
		}

		noMerge, err := cmd.Flags().GetBool("no-merge")
		if err != nil {
			return err
		}
		var ciPipeline map[string]interface{}
		if noMerge == false {
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

		err = os.WriteFile(pipelineFile, data, 0666)
		logger.Info("wrote yaml to file: " + pipelineFile)

		return err
	},
}

func init() {
	Cmd.Flags().String("file", "pipeline.yaml", "output file")
	Cmd.Flags().String("config", "config.kapigen.yaml", "config to use")
	Cmd.Flags().Bool("no-merge", false, "disable dynamic job merge")
	Cmd.Flags().String("mode", version.FILE.Name(), "mode used for versioning: file,gitlab")

}

package pipeline

import (
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/cli"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/config"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/jobs"
	"gitlab.com/kateops/kapigen/cli/internal/version"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gopkg.in/yaml.v3"
)

var GenerateCmd = &cobra.Command{
	Use:              "generate",
	Short:            "Generate pipeline",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.PreparePersistentFlags(cmd)
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
		privateTokenName, err := cmd.Flags().GetString("private-token")
		if err != nil {
			return err
		}
		logger.Info("will create settings")

		logger.Info("will read pipeline config from: " + configPath)
		cmd.SilenceUsage = true
		body, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}

		var pipelineConfig types.PipelineConfig
		err = yaml.Unmarshal(body, &pipelineConfig)
		if err != nil {
			return err
		}
		if privateTokenName == "" {
			privateTokenName = pipelineConfig.PrivateTokenName
		}

		settings := cli.NewSettings(
			cli.SetMode(version.GetModeFromString(mode)),
			cli.SetPrivateToken(privateTokenName),
		)
		pipelineJobs, err := types.LoadJobsFromPipelineConfig(factory.New(settings), pipelineConfig, config.PipelineConfigTypes)
		if err != nil {
			return err
		}
		logger.Info("ci jobs created")

		if pipelineConfig.Noop {
			logger.Info("noop mode activated, will add \"Noop\" job to pipeline")
			pipelineJobs.AddJob(jobs.NewNoop())
		}

		if pipelineConfig.Versioning {
			logger.Info("tag mode activated, will add \"Versioning\" job to pipeline")
			pipelineJobs.AddJob(jobs.NewTag(settings.PrivateToken)).
				AddJob(jobs.NewTagKapigen(settings.PrivateToken))
		}

		noMerge, err := cmd.Flags().GetBool("no-merge")
		if err != nil {
			return err
		}
		if !noMerge {
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
		pipelineJobs.OverwriteTags(pipelineConfig.Tags)
		logger.Info("ci jobs named to be unique")
		return types.JobsToYamLFile(pipelineJobs, pipelineFile)

	},
}

func init() {
	GenerateCmd.Flags().String("file", "pipeline.yaml", "output file")
	GenerateCmd.Flags().String("config", "config.kapigen.yaml", "config to use")
	GenerateCmd.Flags().Bool("no-merge", false, "disable dynamic job merge")
	GenerateCmd.Flags().String("mode", version.FILE.Name(), "mode used for versioning: file,gitlab")

}

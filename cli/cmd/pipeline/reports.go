package pipeline

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	gitlab2 "github.com/xanzy/go-gitlab"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/version"
)

var ReportsCmd = &cobra.Command{
	Use:              "reports",
	Short:            "Get pipeline",
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli.PreparePersistentFlags(cmd)
		logger.Debug("activated verbose mode")

		configPath, err := cmd.Flags().GetString("config")
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
			cli.SetMode(version.GetModeFromString(version.Gitlab.Name())),
			cli.SetPrivateToken(privateTokenName),
		)
		gitlab := factory.New(settings).GetGitlabClient()
		pipelineId, err := strconv.ParseInt(environment.CI_PIPELINE_ID.Get(), 10, 32)
		if err != nil {
			return err
		}
		bridges, res, err := gitlab.Jobs.ListPipelineBridges(environment.CI_PROJECT_ID.Get(), int(pipelineId), nil)
		if res.StatusCode != 200 {
			logger.Error(res.Status)
			return err
		}
		downstreamPipelineIds := []int{}
		for _, bridge := range bridges {
			if bridge.DownstreamPipeline != nil {
				downstreamPipelineIds = append(downstreamPipelineIds, bridge.DownstreamPipeline.ID)
				logger.DebugAny(bridge.DownstreamPipeline)
			}
		}
		var reportJobs wrapper.Array[gitlab2.Job]
		for _, downstreamPipelineId := range downstreamPipelineIds {
			jobs, res, err := gitlab.Jobs.ListPipelineJobs(environment.CI_PROJECT_ID.Get(), downstreamPipelineId, nil)
			if res.StatusCode != 200 {
				logger.Error(res.Status)
				return err
			}
			for _, job := range jobs {
				if job.Artifacts != nil {
					for _, artifact := range job.Artifacts {
						if artifact.FileType == "junit" {
							reportJobs.Push(*job)
							logger.DebugAny(artifact)
						}
					}
				}
			}
		}
		reportJobs.ForEach(func(e *gitlab2.Job) {
			logger.DebugAny(e.Name)
		})
		return err
	},
}

func init() {
	ReportsCmd.Flags().String("config", "config.kapigen.yaml", "config to use")
	ReportsCmd.Flags().String("mode", version.FILE.Name(), "mode used for versioning: file,gitlab")

}

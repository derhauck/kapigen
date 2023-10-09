package pipeline

import (
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/internal/cli"
	"kapigen.kateops.com/internal/gitlab"
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
		configPath := "config.kapigen.yaml"
		body, err := os.ReadFile(configPath)
		if err != nil {
			return err
		}
		var pipelineConfig types.PipelineConfig
		err = yaml.Unmarshal(body, &pipelineConfig)
		if err != nil {
			return err
		}

		var pipelineJobs types.Jobs

		for i := 0; i < len(pipelineConfig.Pipelines); i++ {
			configuration := pipelineConfig.Pipelines[i]
			jobs, err := configuration.Decode(config.PipelineConfigTypes)
			if err != nil {
				return err
			}
			pipelineJobs = append(pipelineJobs, jobs.GetJobs()...)
		}
		logger.Info("pipeline created")
		var ciPipeline = make(map[string]interface{})
		var evaluatedJobs types.Jobs
		var jobsToEvaluate types.Jobs
		jobsToEvaluate = append(jobsToEvaluate, pipelineJobs.GetJobs()...)
		for _, job := range pipelineJobs {
			evaluatedJob, err := job.EvaluateName(&jobsToEvaluate)
			if err != nil {
				return err
			}
			if evaluatedJob != nil {
				evaluatedJobs = append(evaluatedJobs, evaluatedJob)
			} else {
				var resizedJobsToEvaluate types.Jobs
				for i := range jobsToEvaluate {
					if jobsToEvaluate[i] == job && i < len(jobsToEvaluate) {
						var tmp = jobsToEvaluate[i+1:]
						resizedJobsToEvaluate = append(jobsToEvaluate[:i], tmp...)
					}
				}
				jobsToEvaluate = resizedJobsToEvaluate
			}

		}
		for _, job := range evaluatedJobs {
			job.RenderNeeds()
			ciPipeline[job.GetName()] = job.CiJobYaml
		}
		pipeline := gitlab.NewDefaultCiPipeline().Render()
		//ciPipeline["workflow"] = pipeline.Workflow
		ciPipeline["stages"] = pipeline.Stages
		//ciPipeline["default"] = pipeline.Default
		ciPipeline["variables"] = pipeline.Variables
		//logger.DebugAny(ciPipeline)
		data, err := yaml.Marshal(ciPipeline)
		if err != nil {
			return err
		}
		err = os.WriteFile("pipeline.yaml", data, 0777)

		return err
	},
}

func init() {
	Cmd.Flags().Bool("file", false, "output file")
	Cmd.Flags().Bool("config", false, "config to use")

}

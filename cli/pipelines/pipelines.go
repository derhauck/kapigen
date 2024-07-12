package pipelines

import (
	"fmt"
	"os"

	"gitlab.com/kateops/kapigen/cli/internal/pipeline/config"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/pipeline"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gopkg.in/yaml.v3"
)

func ExtendPipelines(pipelines map[types.PipelineType]types.PipelineConfigInterface) {
	for key, pipe := range pipelines {
		logger.Info(fmt.Sprintf("extend pipelines types with: %s", key.String()))
		config.PipelineConfigTypes[key] = pipe
	}
}
func JobsToYamLFile(jobs *types.Jobs, fileName string) error {
	// convert jobs to map
	var ciPipeline = make(map[string]interface{})
	for _, evaluatedJob := range *jobs {
		renderedJob := evaluatedJob.RenderNeeds()
		if renderedJob == nil {
			return fmt.Errorf("job '%s' can not be rendered", evaluatedJob.GetName())
		}
		ciPipeline[evaluatedJob.GetName()] = evaluatedJob.CiJobYaml
	}
	logger.Info("ci job list converted to map")

	// add default pipeline settings
	defaultPipeline := pipeline.NewDefaultCiPipeline().Render()
	ciPipeline["default"] = defaultPipeline.Default
	ciPipeline["workflow"] = defaultPipeline.Workflow
	ciPipeline["stages"] = defaultPipeline.Stages
	ciPipeline["variables"] = defaultPipeline.Variables

	// convert map to yaml
	data, err := yaml.Marshal(ciPipeline)
	if err != nil {
		return err
	}
	logger.Info("converted pipeline to yaml")
	return os.WriteFile(fileName, data, 0666)
}

func CreatePipeline(fn func(jobs *types.Jobs)) {
	jobs := &types.Jobs{}
	fn(jobs)
	evaluatedJobs, err := jobs.EvaluateNames()
	if err != nil {
		logger.ErrorE(err)
		return
	}
	err = JobsToYamLFile(evaluatedJobs, "pipeline.yaml")
	if err != nil {
		logger.ErrorE(err)
		return
	}
}

func ReadPipelineConfig(path string) (*types.PipelineConfig, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pipelineConfig types.PipelineConfig
	err = yaml.Unmarshal(body, &pipelineConfig)
	if err != nil {
		return nil, err
	}
	return &pipelineConfig, nil
}

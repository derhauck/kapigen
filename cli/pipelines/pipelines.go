package pipelines

import (
	"fmt"
	"os"

	"gitlab.com/kateops/kapigen/cli/internal/pipeline/config"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/pipeline"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
	"gopkg.in/yaml.v3"
)

func ExtendPipelinesTypes(pipelines map[types.PipelineType]types.PipelineConfigInterface) {
	for key, pipe := range pipelines {
		logger.Info(fmt.Sprintf("extend pipelines types with: %s", key.String()))
		config.PipelineConfigTypes[key] = pipe
	}
}
func JobsToYamLFile(jobs *types.Jobs, ciPipeline *pipeline.CiPipeline, fileName string) error {
	// convert jobs to map
	var gitlabPipeline = make(map[string]interface{})
	for _, evaluatedJob := range *jobs {
		renderedJob := evaluatedJob.RenderNeeds()
		if renderedJob == nil {
			return fmt.Errorf("job '%s' can not be rendered", evaluatedJob.GetName())
		}
		gitlabPipeline[evaluatedJob.GetName()] = evaluatedJob.CiJobYaml
	}
	logger.Info("ci job list converted to map")

	if ciPipeline == nil {
		return wrapper.DetailedErrorf("ci pipeline can not be nil")
	}

	// add default pipeline settings
	defaultPipeline, err := ciPipeline.Render()
	if err != nil {
		return err
	}
	gitlabPipeline["default"] = defaultPipeline.Default
	gitlabPipeline["workflow"] = defaultPipeline.Workflow
	gitlabPipeline["stages"] = defaultPipeline.Stages
	gitlabPipeline["variables"] = defaultPipeline.Variables

	// convert map to yaml
	data, err := yaml.Marshal(gitlabPipeline)
	if err != nil {
		return err
	}
	logger.Info("converted pipeline to yaml")
	return os.WriteFile(fileName, data, 0666)
}

func CreatePipeline(fn func(jobs *types.Jobs, ciPipeline *pipeline.CiPipeline)) {
	jobs := &types.Jobs{}
	ciPipeline := &pipeline.CiPipeline{}
	fn(jobs, ciPipeline)
	evaluatedJobs, err := jobs.EvaluateNames()
	if err != nil {
		logger.ErrorE(err)
		return
	}
	err = JobsToYamLFile(evaluatedJobs, ciPipeline, "pipeline.yaml")
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

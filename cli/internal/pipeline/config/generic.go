package config

import (
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	generic2 "kapigen.kateops.com/internal/pipeline/jobs/generic"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Generic struct {
	ImageName       string            `yaml:"imageName"`
	Mode            string            `yaml:"mode"`
	Scripts         []string          `yaml:"scripts"`
	Variables       map[string]string `yaml:"variables"`
	Stage           string            `yaml:"stage"`
	Artifacts       job.ArtifactsYaml `yaml:"artifacts"`
	InternalStage   stages.Stage
	InternalChanges []string
}

func (g *Generic) New() types.PipelineConfigInterface {
	return &Generic{}
}

func (g *Generic) Validate() error {

	return nil
}

func (g *Generic) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var allJobs *types.Jobs
	generic, err := generic2.NewGenericJob(g.ImageName, g.InternalStage, g.Scripts)
	if err != nil {
		return nil, err
	}
	return allJobs.AddJob(generic), nil
}

func (g *Generic) Rules() *job.Rules {
	return job.DefaultPipelineRules(g.InternalChanges)
}

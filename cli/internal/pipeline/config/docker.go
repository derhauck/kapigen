package config

import (
	"errors"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Docker struct {
	Path    string `yaml:"path"`
	Context string `yaml:"context"`
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() types.PipelineConfigInterface {
	return d
}

func (d *Docker) LoadPipeline(pipelineType types.PipelineType) types.PipelineBuilderWrapper {

	return types.PipelineBuilderWrapper{
		Builder: DockerBuilder{},
		Config:  d,
		Name:    []string{"Docker"},
		Type:    pipelineType,
	}
}

type DockerBuilder struct {
	Type types.PipelineType
}

func (d DockerBuilder) Build(config types.PipelineConfigInterface) (*gitlab.CiJobs, error) {
	if _, ok := config.(*Docker); ok {
		var tmp gitlab.CiJobs
		return &tmp, nil
	}
	return nil, errors.New("wrong pipeline config")
}

package config

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Infrastructure struct {
	State string `yaml:"state"`
	S3    bool   `yaml:"s3"`
}

func (i *Infrastructure) New() types.PipelineConfigInterface {
	return &Infrastructure{}
}

func (i *Infrastructure) Validate() types.PipelineConfigInterface {
	if i.State == "" {
		i.State = "set-by-validation"
	}
	return i
}

func (i *Infrastructure) LoadPipeline(pipelineType types.PipelineType) types.PipelineBuilderWrapper {
	logger.Debug(fmt.Sprintf("state: %s", i.State))
	logger.Debug(fmt.Sprintf("S3: %s", i.S3))
	return types.PipelineBuilderWrapper{
		Builder: &InfrastructureBuilder{},
		Config:  i,
		Name:    []string{"Infrastructure"},
		Type:    pipelineType,
	}
}

type InfrastructureBuilder struct {
}

func (i *InfrastructureBuilder) Build(config types.PipelineConfigInterface) (*gitlab.CiJobs, error) {
	if _, ok := config.(*Infrastructure); ok {
		var tmp gitlab.CiJobs
		return &tmp, nil
	}

	return nil, errors.New("wrong config type")
}

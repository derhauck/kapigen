package types

import (
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
)

type PipelineType string

func (p PipelineType) String() string {
	return string(p)
}

type ConfigInterface interface {
	Validate() error
}
type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	ConfigInterface
	Build(factory *factory.MainFactory, pipelineType PipelineType, Id string) (*Jobs, error)
	Rules() *job.Rules
}

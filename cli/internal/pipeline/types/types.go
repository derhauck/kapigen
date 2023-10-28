package types

import (
	"kapigen.kateops.com/factory"
)

type PipelineType string

func (p PipelineType) String() string {
	return string(p)
}

type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	Validate() error
	Build(factory *factory.MainFactory, pipelineType PipelineType, Id string) (*Jobs, error)
}

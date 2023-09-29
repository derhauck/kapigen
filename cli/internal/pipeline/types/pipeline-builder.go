package types

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/gitlab"
)

type PipelineBuilderInterface interface {
	Build(config PipelineConfigInterface) (*gitlab.CiJobs, error)
}

type PipelineBuilderWrapper struct {
	Builder PipelineBuilderInterface
	Name    []string
	Config  PipelineConfigInterface
	Type    PipelineType
}

func (p *PipelineBuilderWrapper) Build() (*gitlab.CiJobs, error) {

	if p.Builder == nil {
		return nil, errors.New(fmt.Sprintf("no Pipeline Builder set for type:%s", p.Type))
	}

	return p.Builder.Build(p.Config)
}

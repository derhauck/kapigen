package types

import "gitlab.com/kateops/kapigen/dsl/wrapper"

type PipelineBuilderInterface interface {
	Build(pipelineTypeConfig PipelineTypeConfig) (*Jobs, error)
}

type PipelineBuilderWrapper struct {
	Builder            PipelineBuilderInterface
	PipelineTypeConfig PipelineTypeConfig
}

func (p *PipelineBuilderWrapper) Build() (*Jobs, error) {

	if p.Builder == nil {
		return nil, wrapper.DetailedErrorf("no Pipeline Builder set for type:%s", p.PipelineTypeConfig.GetType())
	}

	return p.Builder.Build(p.PipelineTypeConfig)
}

package types

type PipelineType string

func (p PipelineType) String() string {
	return string(p)
}

type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	Validate() error
	Build(pipelineType PipelineType, Id string) (*Jobs, error)
}

package types

type PipelineType string

type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	Validate() error
	Build(pipelineType PipelineType, Id string) (*Jobs, error)
}

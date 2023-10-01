package types

type PipelineType string

type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	Validate() PipelineConfigInterface
	Build(pipelineType PipelineType, Id string) (*Jobs, error)
}

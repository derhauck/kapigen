package types

type PipelineType string

type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	Validate() error
	Build(path string, pipelineType PipelineType, Id string) (*Jobs, error)
}

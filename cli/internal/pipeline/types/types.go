package types

type PipelineType string

type PipelineConfigInterface interface {
	New() PipelineConfigInterface
	Validate() PipelineConfigInterface
	LoadPipeline(pipelineType PipelineType) PipelineBuilderWrapper //types TODO:
}
type PipelineEntry struct {
	Type   PipelineType `json:"type" yaml:"type"`
	Id     string       `json:"id" yaml:"id"`
	Config interface{}  `json:"config" yaml:"config"`
}

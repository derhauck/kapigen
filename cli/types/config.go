package types

type PipelineTypeConfigInterface interface {
	Decode(configTypes map[PipelineType]PipelineConfigInterface) error
	GetType() PipelineType
}

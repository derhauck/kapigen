package config

import (
	"gitlab.com/kateops/kapigen/cli/types"
)

const (
	DockerPipeline types.PipelineType = "docker"
	PHPPipeline    types.PipelineType = "php"
	GOLANG         types.PipelineType = "golang"
	GENERIC        types.PipelineType = "generic"
)

var PipelineConfigTypes = map[types.PipelineType]types.PipelineConfigInterface{
	DockerPipeline: &Docker{},
	GOLANG:         &Golang{},
	PHPPipeline:    &Php{},
	GENERIC:        &Generic{},
}

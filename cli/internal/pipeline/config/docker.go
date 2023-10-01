package config

import (
	"kapigen.kateops.com/internal/pipeline/types"
)

type Docker struct {
	Path    string `yaml:"path"`
	Context string `yaml:"context"`
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() types.PipelineConfigInterface {
	return d
}

func (i *Docker) Build(pipelineType types.PipelineType, Id string) (*types.Jobs, error) {

	var tmp types.Jobs
	return &tmp, nil

}

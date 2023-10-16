package config

import (
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/golang"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Golang struct {
	ImageName string `yaml:"imageName"`
	Path      string `yaml:"path"`
}

func (g Golang) New() types.PipelineConfigInterface {
	return &Golang{}
}

func (g Golang) Validate() error {

	logger.DebugAny(g.ImageName)
	return nil
}

func (g Golang) Build(pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	test := golang.NewGolangTest(g.ImageName, g.Path)

	return &types.Jobs{
		test,
	}, nil
}

package config

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
)

type PipelineTypeConfig struct {
	Type   types.PipelineType `yaml:"type"`
	Config interface{}        `yaml:"config"`
}

func (p *PipelineTypeConfig) Decode() error {
	logger.Info(fmt.Sprintf("decoding pipeline type %s", p.Type))
	var pipelineConfig = PipelineConfigTypes[p.Type].New()
	if pipelineConfig == nil {
		return errors.New(fmt.Sprintf("no pipeline definition found for type: '%s'", p.Type))
	}
	err := mapstructure.Decode(p.Config, pipelineConfig)
	if err != nil {
		return err
	}
	pipelineConfig.Validate()
	pipelineBuilder := pipelineConfig.LoadPipeline(p.Type)
	jobs, err := pipelineBuilder.Build()
	if err != nil {
		return err
	}

	logger.DebugAny(jobs)

	return nil
}

type PipelineConfig struct {
	Pipelines []PipelineTypeConfig `yaml:"pipelines"`
}

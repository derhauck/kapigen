package types

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"kapigen.kateops.com/internal/logger"
)

type PipelineTypeConfig struct {
	Type       PipelineType `yaml:"type"`
	Config     interface{}  `yaml:"config"`
	PipelineId string       `yaml:"id"`
}

func (p *PipelineTypeConfig) Decode(configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, error) {
	logger.Info(fmt.Sprintf("decoding pipeline type %s", p.Type))
	var pipelineConfig = configTypes[p.Type].New()
	if pipelineConfig == nil {
		return nil, errors.New(fmt.Sprintf("no pipeline definition found for type: '%s'", p.Type))
	}
	err := mapstructure.Decode(p.Config, pipelineConfig)
	if err != nil {
		return nil, err
	}
	pipelineConfig.Validate()
	jobs, err := pipelineConfig.Build(p.Type, p.PipelineId)
	if err != nil {
		return nil, err
	}

	for _, job := range jobs.GetJobs() {
		job.AddName(string(p.Type))
		if p.PipelineId != "" {
			job.AddName(p.PipelineId)
		}
	}

	return jobs, nil
}

func (p *PipelineTypeConfig) GetType() PipelineType {
	return p.Type
}

type PipelineConfig struct {
	Pipelines []PipelineTypeConfig `yaml:"pipelines"`
}

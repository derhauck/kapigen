package types

import (
	"errors"
	"fmt"
	"os"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/logger"
)

type PipelineTypeConfig struct {
	Type       PipelineType `yaml:"type"`
	Config     interface{}  `yaml:"config"`
	PipelineId string       `yaml:"id"`
	Needs      []string     `yaml:"needs"`
	Tags       []string     `yaml:"tags"`
}

func (p *PipelineTypeConfig) Decode(factory *factory.MainFactory, configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, error) {
	logger.Info(fmt.Sprintf("decoding pipeline type %s, id: %s", p.Type, p.PipelineId))
	var pipelineConfig = configTypes[p.Type].New()
	if pipelineConfig == nil {
		return nil, errors.New(fmt.Sprintf("no pipeline definition found for type: '%s'", p.Type))
	}
	err := mapstructure.Decode(p.Config, pipelineConfig)
	if err != nil {
		return nil, err
	}

	jobs, err := GetPipelineJobs(factory, pipelineConfig, p.Type, p.PipelineId)
	if err != nil {
		return nil, err
	}
	for _, job := range jobs.GetJobs() {
		if p.PipelineId != "" {
			job.AddName(p.PipelineId)
		}
		job.AddName(string(p.Type))
		if len(p.Tags) > 0 {
			job.ExternalTags = p.Tags
		}
		err = job.Render()
		if err != nil {
			return nil, err
		}
	}

	return jobs, nil
}

func (p *PipelineTypeConfig) GetType() PipelineType {
	return p.Type
}

type PipelineConfig struct {
	Noop       bool                 `yaml:"noop,omitempty"`
	Versioning bool                 `yaml:"versioning,omitempty"`
	Tags       []string             `yaml:"tags"`
	Pipelines  []PipelineTypeConfig `yaml:"pipelines" yaml:"pipelines"`
}

func GetPipelineJobs(factory *factory.MainFactory, config PipelineConfigInterface, pipelineType PipelineType, pipelineId string) (*Jobs, error) {
	err := config.Validate()
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Pipeline type: %s, id: %s, encountered validation error: %s",
			pipelineType,
			pipelineId,
			err.Error(),
		))
	}

	jobs, err := config.Build(factory, pipelineType, pipelineId)
	if err != nil {
		return nil, errors.New(fmt.Sprintf(
			"Pipeline type: %s, id: %s, encountered build error: %s",
			pipelineType,
			pipelineId,
			err.Error(),
		))
	}
	return jobs.SetPipelineId(pipelineId), nil
}

func LoadJobsFromPipelineConfig(factory *factory.MainFactory, configPath string, configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, *PipelineConfig, error) {
	body, err := os.ReadFile(configPath)

	if err != nil {
		return nil, nil, err
	}

	var pipelineConfig PipelineConfig
	err = yaml.Unmarshal(body, &pipelineConfig)
	if err != nil {
		return nil, nil, err
	}

	jobs, err := pipelineConfig.Decode(factory, configTypes)
	if err != nil {
		return nil, nil, err
	}

	return jobs, &pipelineConfig, nil
}

func (p *PipelineConfig) Decode(factory *factory.MainFactory, configTypes map[PipelineType]PipelineConfigInterface) (*Jobs, error) {

	var pipelineJobs Jobs
	for i := 0; i < len(p.Pipelines); i++ {
		configuration := p.Pipelines[i]
		if configuration.PipelineId == "" {
			value, _ := yaml.Marshal(configuration)
			return nil, errors.New(fmt.Sprintf("no pipeline id set for pipeline %v", string(value)))
		}
		if configuration.Type == "" {
			value, _ := yaml.Marshal(configuration)
			return nil, errors.New(fmt.Sprintf("no pipeline type set for pipeline %v", string(value)))
		}
		jobs, err := configuration.Decode(factory, configTypes)
		if err != nil {
			return nil, err
		}

		for _, pipelineId := range configuration.Needs {
			logger.Info(fmt.Sprintf("PipelineId: '%s' adding need '%s'", configuration.PipelineId, pipelineId))
			if needJobs := pipelineJobs.FindJobsByPipelineId(pipelineId); len(needJobs.GetJobs()) > 0 {
				jobs.SetJobsAsNeed(needJobs)
			}
		}

		pipelineJobs = append(pipelineJobs, jobs.GetJobs()...)
	}

	return &pipelineJobs, nil
}

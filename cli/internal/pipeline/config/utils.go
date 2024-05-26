package config

import (
	"fmt"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/pipeline/types"
)

type SlimDocker struct {
	Path       string            `yaml:"path"`
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile"`
	BuildArgs  map[string]string `yaml:"buildArgs,omitempty"`
}

type Service struct {
	Name      string            `yaml:"name"`
	Variables map[string]string `yaml:"variables"`
	Port      int32             `yaml:"port"`
	ImageName string            `yaml:"imageName"`
	Docker    *SlimDocker       `yaml:"docker"`
}

func (s *Service) Validate() error {

	if s.Name == "" {
		return types.NewMissingArgError("service.name")
	}

	if s.Port <= 0 {
		return types.DetailedErrorf("service: '%s', invalid port %d (must be 1 - 65535)", s.Name, s.Port)
	}

	if s.ImageName == "" && s.Docker == nil {
		return types.NewMissingArgsError("service.imageName", "service.docker")
	}

	return nil
}

func (s *Service) CreateService(factory *factory.MainFactory, Id string, pipelineType types.PipelineType) (*types.Jobs, *job.Service, error) {
	if s.Docker != nil {
		dockerPipeline := &Docker{}
		release := false
		name := fmt.Sprintf("%s-%s", Id, s.Name)
		dockerPipeline.Release = &release
		dockerPipeline.Name = name
		dockerPipeline.Path = s.Docker.Path
		dockerPipeline.Context = s.Docker.Context
		dockerPipeline.Dockerfile = s.Docker.Dockerfile
		dockerPipeline.BuildArgs = s.Docker.BuildArgs
		jobs, err := types.GetPipelineJobs(factory, dockerPipeline, pipelineType, name)
		if err != nil {
			return nil, nil, err
		}
		dockerPipeline.GetFinalImageName()
		service := job.NewService(dockerPipeline.GetFinalImageName(), s.Name, s.Port)

		return jobs, service, nil
	} else {
		service := job.NewService(s.ImageName, s.Name, s.Port)

		return &types.Jobs{}, service, nil
	}
}

type Services []Service

func (s *Services) Validate() error {
	servicePorts := map[int32]bool{}
	for _, service := range *s {

		if err := service.Validate(); err != nil {
			return err
		}
		if servicePorts[service.Port] {
			return types.DetailedErrorf("service: '%s', referencing occupied port: %d", service.Name, service.Port)
		}
		servicePorts[service.Port] = true
	}
	return nil
}

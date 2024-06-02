package config

import (
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

func (s *SlimDocker) DockerConfig() *Docker {
	release := false
	return &Docker{
		Path:       s.Path,
		Context:    s.Context,
		Dockerfile: s.Dockerfile,
		Release:    &release,
		BuildArgs:  s.BuildArgs,
	}
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

func (s *Service) CreateService(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, *job.Service, error) {
	if s.Docker != nil {
		dockerPipeline := s.Docker.DockerConfig()
		jobs, err := types.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, nil, err
		}
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

func (s *Services) AddToJob(factory *factory.MainFactory, pipelineType types.PipelineType, Id string, pipelineJobs *types.Jobs, targetJob *types.Job) error {
	for _, service := range *s {
		jobs, jobService, err := service.CreateService(factory, pipelineType, Id)
		if err != nil {
			return types.DetailedErrorf(err.Error())
		}
		for _, serviceJob := range jobs.GetJobs() {
			targetJob.AddJobAsNeed(serviceJob)
			pipelineJobs.AddJob(serviceJob)
		}
		targetJob.CiJob.Services.Add(jobService)
	}

	return nil
}

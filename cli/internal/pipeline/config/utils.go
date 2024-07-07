package config

import (
	"gitlab.com/kateops/kapigen/cli/factory"
	types2 "gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
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
		return wrapper.NewMissingArgError("service.name")
	}

	if s.Port <= 0 {
		return wrapper.DetailedErrorf("service: '%s', invalid port %d (must be 1 - 65535)", s.Name, s.Port)
	}

	if s.ImageName == "" && s.Docker == nil {
		return wrapper.NewMissingArgsError("service.imageName", "service.docker")
	}

	return nil
}

func (s *Service) CreateService(factory *factory.MainFactory, pipelineType types2.PipelineType, Id string) (*types2.Jobs, *job.Service, error) {
	if s.Docker != nil {
		dockerPipeline := s.Docker.DockerConfig()
		jobs, err := types2.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, nil, err
		}
		service := job.NewService(dockerPipeline.GetFinalImageName(), s.Name, s.Port)

		return jobs, service, nil
	} else {
		service := job.NewService(s.ImageName, s.Name, s.Port)

		return &types2.Jobs{}, service, nil
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
			return wrapper.DetailedErrorf("service: '%s', referencing occupied port: %d", service.Name, service.Port)
		}
		servicePorts[service.Port] = true
	}
	return nil
}

func (s *Services) AddToJob(factory *factory.MainFactory, pipelineType types2.PipelineType, Id string, pipelineJobs *types2.Jobs, targetJob *types2.Job) error {
	for _, service := range *s {
		jobs, jobService, err := service.CreateService(factory, pipelineType, Id)
		if err != nil {
			return wrapper.DetailedErrorf(err.Error())
		}
		for _, serviceJob := range jobs.GetJobs() {
			targetJob.AddJobAsNeed(serviceJob)
			pipelineJobs.AddJob(serviceJob)
		}
		targetJob.CiJob.Services.Add(jobService)
	}

	return nil
}

type JobMode int

const (
	Enabled JobMode = iota
	Permissive
	Disabled
)

var jobModes = map[JobMode]string{
	Enabled:    "enabled",
	Permissive: "permissive",
	Disabled:   "disabled",
}
var JobModeEnum, _ = wrapper.NewEnum[JobMode](jobModes)

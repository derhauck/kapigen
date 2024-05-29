package config

import (
	"fmt"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/php"
	"kapigen.kateops.com/internal/pipeline/types"
)

type PhpComposer struct {
	Path string `yaml:"path"`
	Args string `yaml:"args"`
}

func (p *PhpComposer) Validate() error {
	if p.Path == "" {
		logger.Info("no composer.path set, defaulting to '.'")
		p.Path = "."
	}
	if p.Args == "" {
		logger.Info("no composer.args defaulting to '--no-progress --no-cache --no-interaction'")
		p.Args = "--no-progress --no-cache --no-interaction"
	}
	return nil
}

type Phpunit struct {
	Path string `yaml:"path"`
	Args string `yaml:"args"`
	Bin  string `yaml:"bin"`
}

func (p *Phpunit) Validate(composer *PhpComposer) error {
	if p.Path == "" {
		logger.Info("no phpunit.path set, defaulting to 'composer.path'")
		p.Path = composer.Path
	}
	if p.Args == "" {
		logger.Info("no phpunit.args set")
	}
	if p.Bin == "" {
		logger.Info("no phpunit.bin set, defaulting to '<composer.path>/vendor/bin/phpunit'")
		p.Bin = fmt.Sprintf("%s/vendor/bin/phpunit", composer.Path)
	}
	return nil
}

type Php struct {
	Composer      PhpComposer `yaml:"composer"`
	ImageName     string      `yaml:"ImageName"`
	Phpunit       Phpunit     `yaml:"phpunit"`
	Services      Services    `yaml:"services"`
	Docker        *SlimDocker `yaml:"docker,omitempty"`
	changes       []string
	listenerPorts map[string]int32
}

func (p *Php) New() types.PipelineConfigInterface {
	return &Php{}
}
func (p *Php) Validate() error {
	if err := p.Composer.Validate(); err != nil {
		return err
	}
	if err := p.Phpunit.Validate(&p.Composer); err != nil {
		return err
	}
	if err := p.Services.Validate(); err != nil {
		return types.DetailedErrorE(err)
	}
	p.listenerPorts = make(map[string]int32)
	for _, service := range p.Services {
		p.listenerPorts[service.Name] = service.Port
	}

	if p.Docker != nil {
		if p.Docker.Path == "" {
			return types.NewMissingArgError("docker.path")
		}
		p.ImageName = "docker"
	}

	if p.ImageName == "" && p.Docker == nil {
		return types.NewMissingArgsError("imageName", "docker")
	}
	return nil
}

func (p *Php) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var jobs = &types.Jobs{}
	phpUnitJob, err := php.NewPhpUnit(p.ImageName, p.Composer.Path, p.Composer.Args, p.Phpunit.Path, p.Phpunit.Args, p.Phpunit.Bin, p.listenerPorts)
	p.changes = []string{p.Composer.Path}
	if err != nil {
		return nil, err
	}
	if p.Docker != nil {
		dockerPipeline := &Docker{}
		release := false
		dockerPipeline.Release = &release
		dockerPipeline.Name = Id
		dockerPipeline.Path = p.Docker.Path
		dockerPipeline.Context = p.Docker.Context
		dockerPipeline.Dockerfile = p.Docker.Dockerfile
		dockerPipeline.BuildArgs = p.Docker.BuildArgs
		jobs, err = types.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, err
		}

		for _, currentJob := range jobs.GetJobs() {
			phpUnitJob.AddJobAsNeed(currentJob)
		}
		phpUnitJob.CiJob.Image.Name = dockerPipeline.GetFinalImageName()
		p.changes = append(p.changes, dockerPipeline.Context)
	}
	for _, serviceConfig := range p.Services {
		serviceJobs, service, err := serviceConfig.CreateService(factory, Id, PHPPipeline)
		if err != nil {
			return nil, types.DetailedErrorf(err.Error())
		}

		for _, serviceJob := range serviceJobs.GetJobs() {
			jobs.AddJob(serviceJob)
			phpUnitJob.AddJobAsNeed(serviceJob).
				CiJob.Services.
				Add(service)
		}
	}

	jobs.AddJob(phpUnitJob)
	return jobs, nil
}

func (p *Php) Rules() *job.Rules {
	return &(*job.DefaultPipelineRules(p.changes))
}

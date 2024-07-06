package config

import (
	"fmt"

	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/jobs/php"
	types2 "gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type PhpComposer struct {
	Path string `yaml:"path"`
	Args string `yaml:"args"`
}

func (p *PhpComposer) Validate() error {
	if p.Path == "" {
		logger.Debug("no composer.path set, defaulting to '.'")
		p.Path = "."
	}
	if p.Args == "" {
		logger.Debug("no composer.args defaulting to '--no-progress --no-cache --no-interaction'")
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
		logger.Debug("no phpunit.path set, defaulting to 'composer.path'")
		p.Path = composer.Path
	}
	if p.Args == "" {
		logger.Debug("no phpunit.args set")
	}
	if p.Bin == "" {
		logger.Debug("no phpunit.bin set, defaulting to '<composer.path>/vendor/bin/phpunit'")
		p.Bin = fmt.Sprintf("%s/vendor/bin/phpunit", composer.Path)
	}
	return nil
}

type Php struct {
	Composer              PhpComposer `yaml:"composer"`
	ImageName             string      `yaml:"ImageName"`
	Phpunit               Phpunit     `yaml:"phpunit"`
	Services              Services    `yaml:"services"`
	Docker                *SlimDocker `yaml:"docker,omitempty"`
	InternalChanges       []string
	InternalListenerPorts map[string]int32
}

func (p *Php) New() types2.PipelineConfigInterface {
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
		return wrapper.DetailedErrorE(err)
	}
	p.InternalListenerPorts = make(map[string]int32)
	for _, service := range p.Services {
		p.InternalListenerPorts[service.Name] = service.Port
	}

	if p.Docker != nil {
		if p.Docker.Path == "" {
			return wrapper.NewMissingArgError("docker.path")
		}
		p.ImageName = "docker"
	}

	if p.ImageName == "" && p.Docker == nil {
		return wrapper.NewMissingArgsError("imageName", "docker")
	}
	return nil
}

func (p *Php) Build(factory *factory.MainFactory, pipelineType types2.PipelineType, Id string) (*types2.Jobs, error) {
	var jobs = &types2.Jobs{}
	phpUnitJob, err := php.NewPhpUnit(p.ImageName, p.Composer.Path, p.Composer.Args, p.Phpunit.Path, p.Phpunit.Args, p.Phpunit.Bin, p.InternalListenerPorts)
	p.InternalChanges = []string{p.Composer.Path}
	if err != nil {
		return nil, err
	}
	if p.Docker != nil {
		dockerPipeline := p.Docker.DockerConfig()
		jobs, err = types2.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, err
		}

		for _, currentJob := range jobs.GetJobs() {
			phpUnitJob.AddJobAsNeed(currentJob)
		}
		phpUnitJob.CiJob.Image.Name = dockerPipeline.GetFinalImageName()
		p.InternalChanges = append(p.InternalChanges, dockerPipeline.Context)
	}
	err = p.Services.AddToJob(factory, PHPPipeline, Id, jobs, phpUnitJob)
	if err != nil {
		return nil, err
	}

	jobs.AddJob(phpUnitJob)
	return jobs, nil
}

func (p *Php) Rules() *job.Rules {
	return job.DefaultPipelineRules(p.InternalChanges)
}

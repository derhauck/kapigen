package config

import (
	"errors"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/php"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Php struct {
	ComposerPath   string      `yaml:"composerPath"`
	ImageName      string      `yaml:"ImageName"`
	PhpUnitXmlPath string      `yaml:"phpUnitXmlPath"`
	PhpUnitArgs    string      `yaml:"phpUnitArgs"`
	Docker         *SlimDocker `yaml:"docker,omitempty"`
}

func (p *Php) New() types.PipelineConfigInterface {
	return &Php{}
}
func (p *Php) Validate() error {
	if p.ComposerPath == "" {
		return errors.New("composerPath not set, required")
	}
	if p.PhpUnitXmlPath == "" {
		p.PhpUnitXmlPath = p.ComposerPath
	}
	if p.PhpUnitArgs == "" {
		logger.Info("no phpUnitArgs set")
	}
	if p.Docker != nil && p.Docker.Path == "" {
		return errors.New("docker.path not set, required")
	}
	if p.ImageName == "" && p.Docker == nil {
		return errors.New("imageName and docker not set, one required")
	}
	return nil
}

func (p *Php) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var jobs = &types.Jobs{}
	phpUnitJob, err := php.NewPhpUnit(p.ImageName, p.ComposerPath, p.PhpUnitXmlPath, p.PhpUnitArgs)
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
	}

	jobs.AddJob(phpUnitJob)
	return jobs, nil
}

func (p *Php) Rules() *job.Rules {
	return &(*job.DefaultPipelineRules([]string{p.ComposerPath}))
}

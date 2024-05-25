package config

import (
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/php"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Php struct {
	ComposerPath   string      `yaml:"composerPath"`
	ComposerArgs   string      `yaml:"composerArgs"`
	ImageName      string      `yaml:"ImageName"`
	PhpunitXmlPath string      `yaml:"phpunitXmlPath"`
	PhpunitArgs    string      `yaml:"phpunitArgs"`
	Docker         *SlimDocker `yaml:"docker,omitempty"`
	changes        []string
}

func (p *Php) New() types.PipelineConfigInterface {
	return &Php{}
}
func (p *Php) Validate() error {
	if p.ComposerPath == "" {
		logger.Info("no composerPath set, using default")
		p.ComposerPath = "."
	}
	if p.ComposerArgs == "" {
		logger.Info("no composerArgs set")
	}
	if p.PhpunitXmlPath == "" {
		logger.Info("no phpunitXmlPath set, using same as composerPath")
		p.PhpunitXmlPath = "."
	}
	if p.PhpunitArgs == "" {
		logger.Info("no phpunitArgs set")
	}
	if p.Docker != nil && p.Docker.Path == "" {
		return types.NewMissingArgError("docker.path")
	}
	if p.ImageName == "" && p.Docker == nil {
		return types.NewMissingArgsError("imageName", "docker")
	}
	return nil
}

func (p *Php) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var jobs = &types.Jobs{}
	phpUnitJob, err := php.NewPhpUnit(p.ImageName, p.ComposerPath, p.ComposerArgs, p.PhpunitXmlPath, p.PhpunitArgs)
	p.changes = []string{p.ComposerPath}
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

	jobs.AddJob(phpUnitJob)
	return jobs, nil
}

func (p *Php) Rules() *job.Rules {
	return &(*job.DefaultPipelineRules(p.changes))
}

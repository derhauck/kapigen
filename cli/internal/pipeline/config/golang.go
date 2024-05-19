package config

import (
	"errors"
	"fmt"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/golang"
	"kapigen.kateops.com/internal/pipeline/types"
)

type GolangCoverage struct {
	Packages []string `yaml:"packages"`
}
type GolangDocker struct {
	Path       string `yaml:"path"`
	Context    string `yaml:"context"`
	Dockerfile string `yaml:"dockerfile"`
}

func (g *GolangCoverage) Validate() error {
	if len(g.Packages) == 0 {
		logger.Info("no package declared, using./...")
		g.Packages = []string{"./..."}
	}
	return nil
}

type Golang struct {
	ImageName string          `yaml:"imageName"`
	Path      string          `yaml:"path"`
	Docker    *GolangDocker   `yaml:"docker,omitempty"`
	Coverage  *GolangCoverage `yaml:"coverage,omitempty"`
}

func (g *Golang) New() types.PipelineConfigInterface {
	return &Golang{}
}

func (g *Golang) Validate() error {
	if g.ImageName == "" && g.Docker == nil {
		return errors.New("no imageName or docker config set, required")
	}

	if g.Path == "" {
		logger.Info("no path set, defaulting to '.'")
		g.Path = "."
	}

	if g.Coverage == nil {
		g.Coverage = &GolangCoverage{}
	}

	if err := g.Coverage.Validate(); err != nil {
		return err
	}

	return nil
}

func (g *Golang) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var allJobs = types.Jobs{}
	golangDocker := g.Docker
	docker := &Docker{}
	var test *types.Job
	var err error

	if golangDocker != nil {
		release := false
		docker.Name = Id
		docker.Release = &release
		docker.Name = fmt.Sprintf("golang-%s", Id)
		docker.Path = golangDocker.Path
		docker.Context = golangDocker.Context
		docker.Dockerfile = golangDocker.Dockerfile
	}

	if golangDocker != nil {
		jobs, err := types.GetPipelineJobs(factory, docker, pipelineType, Id)
		if err != nil {
			return nil, err
		}
		test, err = golang.NewUnitTest(docker.GetFinalImageName(), g.Path, g.Coverage.Packages)
		if err != nil {
			return nil, err
		}
		for _, currentJob := range jobs.GetJobs() {
			test.AddJobAsNeed(currentJob)
		}
		allJobs = append(allJobs, jobs.GetJobs()...)
	} else {
		test, err = golang.NewUnitTest(g.ImageName, g.Path, g.Coverage.Packages)
		if err != nil {
			return nil, err
		}
	}

	allJobs = append(allJobs, test)
	return &allJobs, nil
}

func (g *Golang) Rules() *job.Rules {
	return &*job.DefaultPipelineRules(g.Path)
}

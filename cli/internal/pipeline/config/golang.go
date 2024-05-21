package config

import (
	"errors"
	"fmt"
	"os"

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
	Path       string            `yaml:"path"`
	Context    string            `yaml:"context"`
	Dockerfile string            `yaml:"dockerfile"`
	BuildArgs  map[string]string `yaml:"buildArgs,omitempty"`
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
	changes   []string
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

	entries, err := os.ReadDir(g.Path)
	if err != nil {
		logger.ErrorE(err)
		return err
	}
	var isGoMod = false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "go.mod" {
			isGoMod = true
		}
	}
	if isGoMod == false {
		return errors.New("could not find go.mod file in path")
	}

	if g.Coverage == nil {
		g.Coverage = &GolangCoverage{}
	}

	if g.Docker != nil && g.Docker.Path == "" {
		return errors.New("no docker.path set, required")
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
	g.changes = []string{g.Path}
	if golangDocker != nil {
		release := false
		docker.Name = Id
		docker.Release = &release
		docker.Name = fmt.Sprintf("golang-%s", Id)
		docker.Path = golangDocker.Path
		docker.Context = golangDocker.Context
		docker.Dockerfile = golangDocker.Dockerfile
		docker.BuildArgs = golangDocker.BuildArgs
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
		g.changes = append(g.changes, docker.Context)
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
	return &*job.DefaultPipelineRules(g.changes)
}

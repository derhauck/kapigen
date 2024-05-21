package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/golang"
	"kapigen.kateops.com/internal/pipeline/types"
)

type GolangCoverage struct {
	Packages []string `yaml:"packages"`
	Source   string   `yaml:"source"`
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

	if g.Path == "" {
		logger.Info("no path set, defaulting to '.'")
		g.Path = "."
	}
	if g.Coverage == nil {
		g.Coverage = &GolangCoverage{}
	}
	entries, err := os.ReadDir(g.Path)
	var isGoMod = false
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "go.mod" {
			file, err := os.ReadFile(fmt.Sprintf("%s/%s", g.Path, entry.Name()))
			if err != nil {
				return err
			}
			fileString := string(file)

			re := regexp.MustCompile(`go (.*)`)
			match := re.FindStringSubmatch(fileString)
			if len(match) == 0 {
				return fmt.Errorf("go.mod file should include go version")
			}

			logger.DebugAny(match[1])
			g.ImageName = fmt.Sprintf("%s%s:%s", docker.DEPENDENCY_PROXY, "golang", match[1])

			if len(g.Coverage.Packages) == 0 {
				re := regexp.MustCompile(`module (.*)`)
				match := re.FindStringSubmatch(fileString)
				if len(match) == 0 {
					return fmt.Errorf("go.mod file should include module name")
				}
				g.Coverage.Packages = []string{fmt.Sprintf("%s/...", match[1])}
			}
			if g.Coverage.Source == "" {
				g.Coverage.Source = "./..."
			}

			isGoMod = true
		}
	}
	if isGoMod == false {
		return errors.New("could not find go.mod file in path")
	}
	if g.ImageName == "" {
		if err != nil {
			return err
		}

	}

	if g.Docker != nil && g.Docker.Path == "" {
		return errors.New("no docker.path set, required")
	}

	if err := g.Coverage.Validate(); err != nil {
		return err
	}

	if g.ImageName == "" && g.Docker == nil {
		return errors.New("no imageName or docker config set, required")
	}
	return nil
}

func (g *Golang) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var allJobs = types.Jobs{}
	golangDocker := g.Docker
	dockerPipeline := &Docker{}
	var test *types.Job
	var err error
	g.changes = []string{g.Path}
	if golangDocker != nil {
		release := false
		dockerPipeline.Name = Id
		dockerPipeline.Release = &release
		dockerPipeline.Name = fmt.Sprintf("golang-%s", Id)
		dockerPipeline.Path = golangDocker.Path
		dockerPipeline.Context = golangDocker.Context
		dockerPipeline.Dockerfile = golangDocker.Dockerfile
		dockerPipeline.BuildArgs = golangDocker.BuildArgs
		jobs, err := types.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, err
		}
		test, err = golang.NewUnitTest(dockerPipeline.GetFinalImageName(), g.Path, g.Coverage.Packages, g.Coverage.Source)
		if err != nil {
			return nil, err
		}
		for _, currentJob := range jobs.GetJobs() {
			test.AddJobAsNeed(currentJob)
		}
		allJobs = append(allJobs, jobs.GetJobs()...)
		g.changes = append(g.changes, dockerPipeline.Context)
	} else {
		test, err = golang.NewUnitTest(g.ImageName, g.Path, g.Coverage.Packages, g.Coverage.Source)
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

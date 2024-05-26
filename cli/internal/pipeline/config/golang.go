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
	Docker    *SlimDocker     `yaml:"docker"`
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
	if err != nil {
		return err
	}
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

	if g.Docker != nil && g.Docker.Path == "" {
		return types.NewMissingArgError("docker.path")
	}

	if err := g.Coverage.Validate(); err != nil {
		return err
	}

	if g.ImageName == "" && g.Docker == nil {
		return types.NewMissingArgsError("imageName", "docker")
	}
	return nil
}

func (g *Golang) Build(factory *factory.MainFactory, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var allJobs = types.Jobs{}
	dockerPipeline := &Docker{}
	var test *types.Job
	var err error

	test, err = golang.NewUnitTest(g.ImageName, g.Path, g.Coverage.Packages, g.Coverage.Source)
	if err != nil {
		return nil, err
	}
	g.changes = []string{g.Path}
	if g.Docker != nil {
		release := false
		dockerPipeline.Name = Id
		dockerPipeline.Release = &release
		dockerPipeline.Name = fmt.Sprintf("golang-%s", Id)
		dockerPipeline.Path = g.Docker.Path
		dockerPipeline.Context = g.Docker.Context
		dockerPipeline.Dockerfile = g.Docker.Dockerfile
		dockerPipeline.BuildArgs = g.Docker.BuildArgs
		jobs, err := types.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, err
		}
		test.CiJob.SetImageName(dockerPipeline.GetFinalImageName())
		for _, currentJob := range jobs.GetJobs() {
			test.AddJobAsNeed(currentJob)
		}
		allJobs = append(allJobs, jobs.GetJobs()...)
		g.changes = append(g.changes, dockerPipeline.Context)
	}

	allJobs = append(allJobs, test)
	return &allJobs, nil
}

func (g *Golang) Rules() *job.Rules {
	return &*job.DefaultPipelineRules(g.changes)
}

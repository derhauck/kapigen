package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.com/kateops/kapigen/cli/factory"
	"gitlab.com/kateops/kapigen/cli/internal/docker"
	"gitlab.com/kateops/kapigen/cli/internal/pipeline/jobs/golang"
	types2 "gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/environment"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type GolangCoverage struct {
	Packages []string `yaml:"packages"`
	Source   string   `yaml:"source"`
}

func (g *GolangCoverage) Validate() error {
	if len(g.Packages) == 0 {
		logger.Debug("no package declared, using./...")
		g.Packages = []string{"./..."}
	}
	if g.Source == "" {
		logger.Debug("no coverage source declared, using./...")
		g.Source = "./..."
	}
	return nil
}

type GolangLint struct {
	imageName string `yaml:"imageName"`
	Mode      string `yaml:"mode"`
	JobMode   JobMode
}

func (g *GolangLint) Validate() error {
	if g.imageName == "" {
		g.imageName = docker.GOLANG_GOLANGCI_LINT.String()
	}
	if g.Mode == "" {
		logger.Debug("no coverage mode declared, using set")
		g.Mode = JobModeEnum.ValueSafe(Enabled)
	}
	mode, err := JobModeEnum.KeyFromValue(g.Mode)
	if err != nil {
		return wrapper.DetailedErrorE(err)
	}
	g.JobMode = mode

	return nil
}

type Golang struct {
	ImageName string          `yaml:"imageName"`
	Path      string          `yaml:"path"`
	Coverage  *GolangCoverage `yaml:"coverage,omitempty"`
	Lint      *GolangLint     `yaml:"lint,omitempty"`
	Services  Services        `yaml:"services"`
	Docker    *SlimDocker     `yaml:"docker"`
	changes   []string
}

func (g *Golang) New() types2.PipelineConfigInterface {
	return &Golang{}
}

func (g *Golang) Validate() error {

	if g.Path == "" {
		logger.Debug("no path set, defaulting to '.'")
		g.Path = "."
	}
	if g.Coverage == nil {
		g.Coverage = &GolangCoverage{}
	}

	if g.Lint == nil {
		g.Lint = &GolangLint{}
	}

	if g.Docker != nil {
		if g.Docker.Path == "" {
			return wrapper.NewMissingArgError("docker.path")
		}
		g.ImageName = "docker"
	}

	if err := g.Lint.Validate(); err != nil {
		return wrapper.DetailedErrorE(err)
	}
	if err := g.Services.Validate(); err != nil {
		return wrapper.DetailedErrorE(err)
	}
	if err := g.Coverage.Validate(); err != nil {
		return wrapper.DetailedErrorE(err)
	}

	if g.ImageName == "" && g.Docker == nil {
		return wrapper.NewMissingArgsError("imageName", "docker")
	}
	return nil
}

func (g *Golang) Build(factory *factory.MainFactory, pipelineType types2.PipelineType, Id string) (*types2.Jobs, error) {
	var allJobs = types2.Jobs{}

	golangUnitTestJob, err := golang.NewUnitTest(g.ImageName, g.Path, g.Coverage.Packages, g.Coverage.Source)
	golangLint := golang.Lint(g.Lint.imageName, g.Path)
	if err != nil {
		return nil, err
	}
	g.changes = []string{g.Path}
	if g.Docker != nil {
		dockerPipeline := g.Docker.DockerConfig()
		jobs, err := types2.GetPipelineJobs(factory, dockerPipeline, pipelineType, Id)
		if err != nil {
			return nil, err
		}
		golangUnitTestJob.CiJob.SetImageName(dockerPipeline.GetFinalImageName())
		for _, currentJob := range jobs.GetJobs() {
			golangUnitTestJob.AddJobAsNeed(currentJob)
		}
		allJobs = append(allJobs, jobs.GetJobs()...)
		g.changes = append(g.changes, dockerPipeline.Context)
	}
	if g.Lint.JobMode == Permissive {
		golangLint.CiJob.AllowFailure.AllowAll()
	}

	if g.Lint.JobMode != Disabled {
		allJobs.AddJob(golangLint)
	}

	err = g.Services.AddToJob(factory, PHPPipeline, Id, &allJobs, golangUnitTestJob)
	if err != nil {
		return nil, err
	}
	return allJobs.
		AddJob(golangUnitTestJob), nil
}

func (g *Golang) Rules() *job.Rules {
	return job.DefaultPipelineRules(g.changes)
}

func GolangAutoConfig() *Golang {
	config := &Golang{}
	files := SearchPath(environment.CI_PROJECT_DIR.Get(), "go.mod", []string{})
	logger.DebugAny(files)
	for _, fileName := range files {
		dir, _ := filepath.Split(fileName)
		dir, found := strings.CutPrefix(dir, fmt.Sprintf("%s/", environment.CI_PROJECT_DIR.Get()))
		if !found {
			return nil
		}
		dir, _ = strings.CutSuffix(dir, "/")
		if dir == "" {
			dir = "."
		}
		config.Path = dir
		file, err := os.ReadFile(fileName)
		if err != nil {
			return nil
		}
		fileString := string(file)

		re := regexp.MustCompile(`go (.*)`)
		match := re.FindStringSubmatch(fileString)
		if len(match) == 0 {
			return nil // fmt.Errorf("go.mod file should include go version")
		}
		fmt.Println(match[1])
		config.ImageName = fmt.Sprintf("%s%s:%s", docker.DEPENDENCY_PROXY, "golang", match[1])
	}
	return config
}

func SearchPath(path string, name string, entries []string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println(err.Error())
		return entries
	}
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if file.IsDir() {
			entries = SearchPath(fmt.Sprintf("%s/%s", path, file.Name()), name, entries)
			continue
		}
		if file.Name() == name {
			entries = append(entries, fmt.Sprintf("%s/%s", path, file.Name()))
		}
	}

	return entries
}

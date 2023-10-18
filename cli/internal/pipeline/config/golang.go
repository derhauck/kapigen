package config

import (
	"kapigen.kateops.com/internal/pipeline/jobs/golang"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Golang struct {
	ImageName string  `yaml:"imageName"`
	Path      string  `yaml:"path"`
	Docker    *Docker `yaml:"docker,omitempty"`
}

func (g Golang) New() types.PipelineConfigInterface {
	return &Golang{}
}

func (g Golang) Validate() error {

	return nil
}

func (g Golang) Build(pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var allJobs = types.Jobs{}
	test := golang.NewGolangTest(g.ImageName, g.Path)
	docker := g.Docker
	if docker != nil {
		jobs, err := types.GetPipelineJobs(docker, pipelineType, Id)
		if err != nil {
			return nil, err
		}
		for _, job := range jobs.GetJobs() {
			test.AddNeed(job)
		}
		allJobs = append(allJobs, jobs.GetJobs()...)

	}
	allJobs = append(allJobs, test)
	return &allJobs, nil
}

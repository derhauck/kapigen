package config

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/pipeline/jobs/docker"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Docker struct {
	Path       string `yaml:"path"`
	Context    string `yaml:"context"`
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() error {
	if d.Path == "" {
		return errors.New("need path set")
	}

	if d.Dockerfile == "" {
		d.Dockerfile = "Dockerfile"
	}

	if d.Context == "" {
		d.Context = d.Path
	}

	if d.Name == "" {

	}

	return nil
}

func (d *Docker) Build(pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var jobs = types.Jobs{docker.NewBuildkitBuild(
		d.Path,
		d.Context,
		d.Dockerfile,
		d.DefaultRegistry(),
	)}
	return &jobs, nil

}

func (d *Docker) DefaultRegistry() string {
	if d.Name != "" {
		return fmt.Sprintf("${CI_REGISTRY_IMAGE}/%s", d.Name)
	}
	return "${CI_REGISTRY_IMAGE}"

}

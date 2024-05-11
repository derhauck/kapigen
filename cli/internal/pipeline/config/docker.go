package config

import (
	"errors"
	"fmt"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/docker"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Docker struct {
	Path       string `yaml:"path"`
	Context    string `yaml:"context"`
	Name       string `yaml:"name"`
	Dockerfile string `yaml:"dockerfile"`
	Release    *bool  `yaml:"release,omitempty"`
	ImageName  string
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() error {
	if d.Path == "" {
		return errors.New("no path set, required")
	}

	if d.Dockerfile == "" {
		d.Dockerfile = "Dockerfile"
	}

	if d.Context == "" {
		logger.Info("no context set, using path")
		d.Context = d.Path
	}

	if d.Name == "" {
		logger.Info("no name set, using container root registry")
	}

	if d.Release == nil {
		logger.Info("no release set, defaulting to false")
		tmp := false
		d.Release = &tmp
	}

	return nil
}

func (d *Docker) Build(factory *factory.MainFactory, _ types.PipelineType, _ string) (*types.Jobs, error) {
	controller := factory.GetVersionController()
	tag := controller.GetCurrentPipelineTag(d.Path)
	var destination []string
	destination = append(destination, d.DefaultRegistry(tag))
	d.ImageName = d.DefaultRegistry(tag)
	if environment.IsRelease() && *d.Release {
		destination = append(destination, d.DefaultRegistry("latest"))
	}

	build := docker.NewDaemonlessBuildkitBuild(
		d.Path,
		d.Context,
		d.Dockerfile,
		destination,
	)
	return &types.Jobs{build}, nil
}

func (d *Docker) GetFinalImageName() string {
	return d.ImageName
}
func (d *Docker) DefaultRegistry(tag string) string {
	if tag == "" {
		tag = "latest"
	}
	if d.Name != "" {
		return fmt.Sprintf("${CI_REGISTRY_IMAGE}/%s:%s", d.Name, tag)
	}
	return fmt.Sprintf("${CI_REGISTRY_IMAGE}:%s", tag)
}

package config

import (
	"fmt"

	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/jobs/docker"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Docker struct {
	Path          string            `yaml:"path"`
	Context       string            `yaml:"context"`
	Name          string            `yaml:"name"`
	Dockerfile    string            `yaml:"dockerfile"`
	Release       *bool             `yaml:"release,omitempty"`
	BuildArgs     map[string]string `yaml:"buildArgs,omitempty"`
	ImageName     string            `yaml:"imageName"`
	PushImageName string
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() error {
	if d.Path == "" {
		logger.Info("no path set, defaulting to '.'")
		d.Path = "."
	}

	if d.Dockerfile == "" {
		d.Dockerfile = "Dockerfile"
	}

	if d.Context == "" {
		logger.Info("no context set, defaulting to path")
		d.Context = d.Path
	}

	if d.Name == "" {
		logger.Info("no name set, defaulting to container root registry")
	}

	if d.Release == nil {
		logger.Info("no release set, defaulting to true")
		tmp := true
		d.Release = &tmp
	}

	return nil
}

func (d *Docker) Build(factory *factory.MainFactory, _ types.PipelineType, _ string) (*types.Jobs, error) {
	controller := factory.GetVersionController()
	tag := controller.GetCurrentPipelineTag(d.Path)
	var destination []string
	destination = append(destination, d.DefaultRegistry(tag))
	d.PushImageName = d.DefaultRegistry(tag)
	if environment.IsRelease() {
		destination = append(destination, d.DefaultRegistry("latest"))
	}
	buildargs := []string{}
	for key, value := range d.BuildArgs {
		buildargs = append(buildargs, fmt.Sprintf("%s=\"%s\"", key, value))
	}

	build := docker.NewDaemonlessBuildkitBuild(
		d.ImageName,
		d.Path,
		d.Context,
		d.Dockerfile,
		destination,
		buildargs,
	)
	return &types.Jobs{build}, nil
}

func (d *Docker) GetFinalImageName() string {
	return d.PushImageName
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

func (d *Docker) Rules() *job.Rules {
	rules := &job.Rules{}
	if *d.Release {
		rules.AddRules(*job.DefaultOnlyReleasePipelineRules())
	}
	rules.AddRules(*job.DefaultMergeRequestRules(d.Context))
	return rules
}

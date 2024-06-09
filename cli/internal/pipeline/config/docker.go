package config

import (
	"fmt"
	"hash/fnv"

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
	Id            string
}

func (d *Docker) New() types.PipelineConfigInterface {
	return &Docker{}
}

func (d *Docker) Validate() error {
	if d.Path == "" {
		logger.Debug("no path set, defaulting to '.'")
		d.Path = "."
	}

	if d.Dockerfile == "" {
		d.Dockerfile = "Dockerfile"
	}

	if d.Context == "" {
		logger.Debug("no context set, defaulting to path")
		d.Context = d.Path
	}

	if d.Name == "" {
		logger.Debug("no name set, defaulting to container root registry")
	}

	if d.Release == nil {
		logger.Debug("no release set, defaulting to true")
		tmp := true
		d.Release = &tmp
	}

	if *d.Release == false {
		configRepresentation := fmt.Sprintf("%s-%s-%s-%s", d.Path, d.Context, d.Dockerfile, d.BuildArgs)
		hasher := fnv.New32a()
		_, err := hasher.Write([]byte(configRepresentation))
		if err != nil {
			return err
		}
		d.Name = fmt.Sprintf("%d", hasher.Sum32())
	}

	return nil
}

func (d *Docker) Build(factory *factory.MainFactory, _ types.PipelineType, Id string) (*types.Jobs, error) {
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

	build, err := docker.NewDaemonlessBuildkitBuild(
		d.ImageName,
		d.Path,
		d.Context,
		d.Dockerfile,
		destination,
		buildargs,
	)
	if err != nil {
		return nil, err
	}
	build.AddName(Id)
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
	rules.AddRules(*job.DefaultMergeRequestRules([]string{d.Context}))
	return rules
}

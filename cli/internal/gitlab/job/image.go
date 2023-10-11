package job

import (
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/images"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type Image struct {
	Name       docker.Image
	Entrypoint wrapper.StringSlice
	PullPolicy images.PullPolicy
}

func (i *Image) GetRenderedValue() *ImageYaml {
	return NewImageYaml(i)
}

type ImageYaml struct {
	Name       string   `yaml:"name"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
	PullPolicy string   `yaml:"pull_policy"`
}

func NewImageYaml(ci *Image) *ImageYaml {
	return &ImageYaml{
		Name:       ci.Name.Image(),
		Entrypoint: ci.Entrypoint.Get(),
		PullPolicy: ci.PullPolicy.String(),
	}
}

package gitlab

import (
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type ImagePullPolicy int

const (
	ImagePullPolicyAlways ImagePullPolicy = iota
	ImagePullPolicyIfNotPresent
	ImagePullPolicyNever
)

func (c ImagePullPolicy) ImagePullPolicy() string {
	return []string{
		"always",
		"if-not-present",
		"never",
	}[c]
}

type Image struct {
	Name       docker.Image
	Entrypoint wrapper.StringSlice
	PullPolicy ImagePullPolicy
}

func (i *Image) GetRenderedValue() *ImageYaml {
	return NewImageYaml(i)
}

type ImageYaml struct {
	Name       string   `yaml:"name"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
	PullPolicy string   `yaml:"pull_policy"`
}

func NewImageYaml(image *Image) *ImageYaml {
	return &ImageYaml{
		Name:       image.Name.Image(),
		Entrypoint: image.Entrypoint.Get(),
		PullPolicy: image.PullPolicy.ImagePullPolicy(),
	}
}

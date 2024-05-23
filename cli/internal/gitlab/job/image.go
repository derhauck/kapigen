package job

import (
	"errors"

	"kapigen.kateops.com/internal/gitlab/images"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type Image struct {
	Name       string
	Entrypoint wrapper.StringSlice
	PullPolicy images.PullPolicy
}

//func NewImage() *Image {
//	return &Image{
//		Name:       "",
//		Entrypoint: *wrapper.NewStringSlice(),
//		PullPolicy: images.Always,
//	}
//}

func (i *Image) GetRenderedValue() (*ImageYaml, error) {
	return NewImageYaml(i)
}

type ImageYaml struct {
	Name       string   `yaml:"name"`
	Entrypoint []string `yaml:"entrypoint,omitempty"`
	PullPolicy string   `yaml:"pull_policy"`
}

func NewImageYaml(ci *Image) (*ImageYaml, error) {
	if ci.Name == "" {
		return nil, errors.New("image name must be set")
	}

	return &ImageYaml{
		Name:       ci.Name,
		Entrypoint: ci.Entrypoint.Get(),
		PullPolicy: ci.PullPolicy.String(),
	}, nil
}

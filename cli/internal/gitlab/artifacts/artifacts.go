package artifacts

import (
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type Ci struct {
	Paths     wrapper.StringSlice
	exclude   wrapper.StringSlice
	ExpireIn  string
	ExposeAs  string
	Name      string
	Reports   interface{}
	Untracked bool
	When      rules.WhenWrapper
}

func (a *Ci) Render() *Yaml {
	return NewArtifactsYaml(a)
}

type Yaml struct {
	Paths     []string    `yaml:"paths"`
	exclude   []string    `yaml:"exclude,omitempty"`
	ExpireIn  *string     `yaml:"expire_in,omitempty"`
	ExposeAs  *string     `yaml:"expose_as,omitempty"`
	Name      *string     `yaml:"name,omitempty"`
	Reports   interface{} `yaml:"reports,omitempty"`
	Untracked bool        `yaml:"untracked"`
	When      string      `yaml:"when"`
}

func NewArtifactsYaml(artifacts *Ci) *Yaml {
	return &Yaml{
		Paths:     artifacts.Paths.Get(),
		exclude:   artifacts.exclude.Get(),
		ExpireIn:  &artifacts.ExpireIn,
		ExposeAs:  &artifacts.ExposeAs,
		Name:      &artifacts.Name,
		Reports:   artifacts.Reports,
		Untracked: artifacts.Untracked,
		When:      artifacts.When.Get(),
	}
}

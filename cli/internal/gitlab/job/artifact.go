package job

import (
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type Artifact struct {
	Paths     wrapper.StringSlice
	Exclude   wrapper.StringSlice
	ExpireIn  string
	ExposeAs  string
	Name      string
	Reports   artifact.Reports
	Untracked bool
	When      WhenWrapper
}

func NewArtifact() Artifact {
	return Artifact{}
}
func (a *Artifact) validate() (bool, error) {
	if len(a.Paths.Get()) == 0 {
		return false, nil
	}

	return true, nil
}

func (a *Artifact) Render() (*ArtifactYaml, error) {
	if ok, err := a.validate(); !ok {
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
	return NewArtifactsYaml(a), nil
}

type ArtifactYaml struct {
	Paths     []string              `yaml:"paths"`
	Exclude   []string              `yaml:"exclude,omitempty"`
	ExpireIn  string                `yaml:"expire_in,omitempty"`
	ExposeAs  string                `yaml:"expose_as,omitempty"`
	Name      string                `yaml:"name,omitempty"`
	Reports   *artifact.ReportsYaml `yaml:"reports,omitempty"`
	Untracked bool                  `yaml:"untracked"`
	When      string                `yaml:"when"`
}

func NewArtifactsYaml(artifacts *Artifact) *ArtifactYaml {

	return &ArtifactYaml{
		Paths:     artifacts.Paths.Get(),
		Exclude:   artifacts.Exclude.Get(),
		ExpireIn:  artifacts.ExpireIn,
		ExposeAs:  artifacts.ExposeAs,
		Name:      artifacts.Name,
		Reports:   artifacts.Reports.Render(),
		Untracked: artifacts.Untracked,
		When:      artifacts.When.Get(),
	}
}

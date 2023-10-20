package job

import (
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
)

type CiJob struct {
	Artifact     Artifact          `yaml:"artifacts"`
	AfterScript  AfterScript       `yaml:"after_script"`
	BeforeScript BeforeScript      `yaml:"before_script"`
	Script       Script            `yaml:"script"`
	AllowFailure AllowFailure      `yaml:"allow_failure"`
	Cache        Cache             `yaml:"cache"`
	Variables    map[string]string `yaml:"variables"`
	Tags         Tags              `yaml:"tags"`
	Image        Image
	Rules        Rules
	Stage        stages.Stage `yaml:"stage"`
	Services     Services     `yaml:"services"`
	Coverage     string       `yaml:"coverage"`
}

func (c *CiJob) Render(needs *NeedsYaml) (*CiJobYaml, error) {
	return NewCiJobYaml(c, needs)
}

type CiJobs []*CiJob

type CiJobYaml struct {
	Artifact     *ArtifactYaml     `yaml:"artifacts,omitempty" json:"artifacts,omitempty"`
	AfterScript  []string          `yaml:"after_script,omitempty" json:"after_script,omitempty"`
	AllowFailure any               `yaml:"allow_failure,omitempty" json:"allow_failure,omitempty"`
	BeforeScript []string          `yaml:"before_script,omitempty" json:"before_script,omitempty"`
	Cache        *CacheYaml        `yaml:"cache,omitempty" json:"cache,omitempty"`
	Script       []string          `yaml:"script" json:"script"`
	Needs        *NeedsYaml        `yaml:"needs" json:"needs"`
	Variables    map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`
	Image        *ImageYaml        `yaml:"image" json:"image"`
	Rules        *RulesYaml        `yaml:"rules" json:"rules"`
	Stage        string            `yaml:"stage" json:"stage"`
	Services     *ServiceYamls     `yaml:"services,omitempty" json:"services,omitempty"`
	Tags         []string          `yaml:"tags" json:"tags"`
	Coverage     string            `yaml:"coverage,omitempty" json:"coverage"`
}

func (c *CiJobYaml) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		logger.ErrorE(err)
		return ""
	}
	return string(data)
}
func NewCiJobYaml(job *CiJob, needs *NeedsYaml) (*CiJobYaml, error) {
	var err error
	artifact, err := job.Artifact.Render()
	if err != nil {
		return nil, err
	}
	cache, err := job.Cache.GetRenderedValue()
	if err != nil {
		return nil, err
	}

	stage := job.Stage
	if stage < stages.DYNAMIC {
		stage = stages.DYNAMIC
	}

	return &CiJobYaml{
		Artifact:     artifact,
		AfterScript:  job.AfterScript.GetRenderedValue(),
		AllowFailure: job.AllowFailure.Get(),
		BeforeScript: job.BeforeScript.GetRenderedValue(),
		Cache:        cache,
		Script:       job.Script.GetRenderedValue(),
		Needs:        needs,
		Variables:    job.Variables,
		Image:        job.Image.GetRenderedValue(),
		Rules:        job.Rules.GetRenderedValue(),
		Stage:        stage.String(),
		Services:     job.Services.Render(),
		Tags:         job.Tags.Render(),
		Coverage:     job.Coverage,
	}, nil
}

type NeedYaml struct {
	Optional bool   `yaml:"optional"`
	Job      string `yaml:"job"`
}

type NeedsYaml []*NeedYaml

func (n *NeedsYaml) GetNeeds() []*NeedYaml {
	return *n
}

func NewNeedsYaml(needs []*NeedYaml) *NeedsYaml {
	var newNeeds NeedsYaml
	for _, need := range needs {
		newNeeds = append(newNeeds, need)
	}
	return &newNeeds
}

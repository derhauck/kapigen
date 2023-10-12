package gitlab

import (
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/images"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
)

type CiJob struct {
	AfterScript  job.AfterScript   `yaml:"after_script"`
	BeforeScript job.BeforeScript  `yaml:"before_script"`
	Script       job.Script        `yaml:"script"`
	AllowFailure job.AllowFailure  `yaml:"allow_failure"`
	Cache        job.Cache         `yaml:"cache"`
	Variables    map[string]string `yaml:"variables"`
	Tags         job.Tags          `yaml:"tags"`
	Image        job.Image
	Rules        job.Rules
	Stage        stages.Stage `yaml:"stage"`
	Services     job.Services `yaml:"services"`
}

func (c *CiJob) AddVariable(key string, value string) *CiJob {
	if c.Variables == nil {
		c.Variables = map[string]string{}
	}
	c.Variables[key] = value
	return c
}

func (c *CiJob) AddService(service *job.Service) {
	c.Services.Add(service)
}

func NewCiJob(imageName docker.Image) *CiJob {
	return &CiJob{
		Script: job.NewScript(),
		Cache:  job.NewCache(),
		Image: job.Image{
			Name:       imageName,
			PullPolicy: images.Always,
		},
		Stage:        stages.NewStage(),
		AfterScript:  job.NewAfterScript(),
		BeforeScript: job.NewBeforeScript(),
	}
}

func (c *CiJob) Render(needs *NeedsYaml) *CiJobYaml {
	return NewCiJobYaml(c, needs)
}

type CiJobs []*CiJob

type CiJobYaml struct {
	AfterScript  []string          `yaml:"after_script,omitempty" json:"after_script,omitempty"`
	AllowFailure any               `yaml:"allow_failure,omitempty" json:"allow_failure,omitempty"`
	BeforeScript []string          `yaml:"before_script,omitempty" json:"before_script,omitempty"`
	Cache        *job.CacheYaml    `yaml:"cache,omitempty" json:"cache,omitempty"`
	Script       []string          `yaml:"script" json:"script"`
	Needs        *NeedsYaml        `yaml:"needs" json:"needs"`
	Variables    map[string]string `yaml:"variables,omitempty" json:"variables,omitempty"`
	Image        *job.ImageYaml    `yaml:"image" json:"image"`
	Rules        *job.RulesYaml    `yaml:"rules" json:"rules"`
	Stage        string            `yaml:"stage" json:"stage"`
	Services     *job.ServiceYamls `yaml:"services,omitempty" json:"services,omitempty"`
	Tags         []string          `yaml:"tags" json:"tags"`
}

func (c *CiJobYaml) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		logger.ErrorE(err)
		return ""
	}
	return string(data)
}
func NewCiJobYaml(job *CiJob, needs *NeedsYaml) *CiJobYaml {
	return &CiJobYaml{
		AfterScript:  job.AfterScript.GetRenderedValue(),
		AllowFailure: job.AllowFailure.Get(),
		BeforeScript: job.BeforeScript.GetRenderedValue(),
		Cache:        job.Cache.GetRenderedValue(),
		Script:       job.Script.GetRenderedValue(),
		Needs:        needs,
		Variables:    job.Variables,
		Image:        job.Image.GetRenderedValue(),
		Rules:        job.Rules.GetRenderedValue(),
		Stage:        job.Stage.String(),
		Services:     job.Services.Render(),
		Tags:         job.Tags.Render(),
	}
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

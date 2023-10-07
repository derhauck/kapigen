package gitlab

import (
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/gitlab/services"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
)

type CiJob struct {
	AfterScript  AfterScript       `yaml:"after_script"`
	BeforeScript BeforeScript      `yaml:"before_script"`
	Script       Script            `yaml:"script"`
	AllowFailure AllowFailure      `yaml:"allow_failure"`
	Cache        Cache             `yaml:"cache"`
	Variables    map[string]string `yaml:"variables"`
	Image        Image
	Rules        rules.Rules
	Stage        stages.Stage `yaml:"stage"`
	Services     services.Cis `yaml:"services"`
}

func (c *CiJob) AddVariable(key string, value string) *CiJob {
	if c.Variables == nil {
		c.Variables = map[string]string{}
	}
	c.Variables[key] = value
	return c
}

func (c *CiJob) AddService(service *services.Ci) {
	c.Services.Add(service)
}

func NewCiJob(imageName docker.Image) *CiJob {
	return &CiJob{
		Script: NewScript(),
		Cache:  NewCache(nil),
		Image: Image{
			Name:       imageName,
			PullPolicy: ImagePullPolicyAlways,
		},
		Stage: stages.NewStage(),
	}
}

func (c *CiJob) Render(needs *NeedsYaml) *CiJobYaml {
	return NewCiJobYaml(c, needs)
}

type CiJobs []*CiJob

type CiJobYaml struct {
	AfterScript  []string          `yaml:"after_script"`
	AllowFailure any               `yaml:"allow_failure,omitempty"`
	BeforeScript []string          `yaml:"before_script"`
	Cache        *CacheYaml        `yaml:"cache,omitempty"`
	Script       []string          `yaml:"script"`
	Needs        *NeedsYaml        `yaml:"needs"`
	Variables    map[string]string `yaml:"variables"`
	Image        *ImageYaml        `yaml:"image"`
	Rules        *rules.RulesYaml  `yaml:"rules"`
	Stage        string            `yaml:"stage"`
	Services     *services.Yamls   `yaml:"services,omitempty"`
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
		AfterScript:  job.AfterScript.getRenderedValue(),
		AllowFailure: job.AllowFailure.Get(),
		BeforeScript: job.BeforeScript.Value.Get(),
		Cache:        job.Cache.getRenderedValue(),
		Script:       job.Script.getRenderedValue(),
		Needs:        needs,
		Variables:    job.Variables,
		Image:        job.Image.GetRenderedValue(),
		Rules:        job.Rules.GetRenderedValue(),
		Stage:        job.Stage.Name(),
		Services:     job.Services.Render(),
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

package gitlab

import (
	"gopkg.in/yaml.v3"
	"kapigen.kateops.com/internal/logger"
)

type CiJob struct {
	AfterScript  AfterScript  `yaml:"after_script"`
	BeforeScript BeforeScript `yaml:"before_script"`
	Script       Script       `yaml:"script"`
	AllowFailure AllowFailure `yaml:"allow_failure"`
	Cache        Cache        `yaml:"cache"`
}

func NewCiJob() *CiJob {
	return &CiJob{
		Script: NewScript(),
		Cache:  Cache{},
	}
}

func (c *CiJob) Render() *CiJobYaml {
	return NewCiJobYaml(c)
}

type CiJobs []*CiJob

type CiJobYaml struct {
	AfterScript  []string   `yaml:"after_script"`
	AllowFailure []int32    `yaml:"allow_failure"`
	BeforeScript []string   `yaml:"before_script"`
	Cache        *CacheYaml `yaml:"cache"`
	Script       []string   `yaml:"script"`
}

func (c *CiJobYaml) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		logger.ErrorE(err)
		return ""
	}
	return string(data)
}
func NewCiJobYaml(job *CiJob) *CiJobYaml {
	return &CiJobYaml{
		AfterScript:  job.AfterScript.getRenderedValue(),
		AllowFailure: job.AllowFailure.Get(),
		BeforeScript: job.BeforeScript.Value.Get(),
		Cache:        job.Cache.getRenderedValue(),
		Script:       job.Script.getRenderedValue(),
	}
}

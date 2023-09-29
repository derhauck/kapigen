package gitlab

type CiJob struct {
	AfterScript  *AfterScript  `yaml:"after_script"`
	BeforeScript *BeforeScript `yaml:"before_script"`
	Script       *Script       `yaml:"script"`
	AllowFailure *AllowFailure `yaml:"allow_failure"`
	Cache        *Cache        `yaml:"cache"`
}

func (c *CiJob) render() *CiJobYaml {
	return NewCiJobYaml(c)
}

type CiJobs []CiJob

type CiJobYaml struct {
	AfterScript  []string  `yaml:"after_script"`
	AllowFailure []int32   `yaml:"allow_failure"`
	BeforeScript []string  `yaml:"before_script"`
	Cache        CacheYaml `yaml:"cache"`
}

func NewCiJobYaml(job *CiJob) *CiJobYaml {
	return &CiJobYaml{
		AfterScript:  job.AfterScript.getRenderedValue(),
		AllowFailure: job.AllowFailure.Get(),
		BeforeScript: job.BeforeScript.Value.Get(),
		Cache:        job.Cache.getRenderedValue(),
	}
}

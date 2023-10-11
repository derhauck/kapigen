package job

import (
	"fmt"
	"kapigen.kateops.com/internal/gitlab/cache"
	"kapigen.kateops.com/internal/gitlab/environment"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type CacheYaml struct {
	Key       string   `yaml:"key"`
	Paths     []string `yaml:"paths"`
	Unprotect bool     `yaml:"unprotect"`
	Policy    string   `yaml:"policy"`
}

type Cache struct {
	Key       string              `yaml:"key"`
	Paths     wrapper.StringSlice `yaml:"paths"`
	Unprotect bool                `yaml:"unprotect"`
	Policy    cache.Policy        `yaml:"policy"`
	Active    bool
}

func (c *Cache) SetActive() *Cache {
	c.Active = true
	return c
}
func (c *Cache) SetPolicy(policy cache.Policy) *Cache {
	c.Policy = policy
	return c
}

func (c *Cache) SetDefaultCacheKey(path string, pipelineType string) {
	c.Key = fmt.Sprintf("%s-%s-%s", environment.Get(environment.CI_MERGE_REQUEST_ID), path, pipelineType)
}
func NewCache() Cache {
	return Cache{
		Paths: wrapper.StringSlice{
			Value: nil,
		},
		Policy:    cache.PullPush,
		Unprotect: true,
	}
}

func (c *Cache) GetRenderedValue() *CacheYaml {
	if c.Key == "" || c.Active == false {
		return nil
	}
	return &CacheYaml{
		Key:       c.Key,
		Paths:     c.Paths.Get(),
		Unprotect: c.Unprotect,
		Policy:    c.Policy.String(),
	}
}

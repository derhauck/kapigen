package job

import (
	"errors"
	"fmt"

	"gitlab.com/kateops/kapigen/dsl/environment"
	"gitlab.com/kateops/kapigen/dsl/gitlab/cache"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type CacheYaml struct {
	Key       string   `yaml:"key"`
	Paths     []string `yaml:"paths"`
	Unprotect bool     `yaml:"unprotect"`
	Policy    string   `yaml:"policy"`
}

type Cache struct {
	Key       string                `yaml:"key"`
	Paths     wrapper.Array[string] `yaml:"paths"`
	Unprotect bool                  `yaml:"unprotect"`
	Policy    cache.Policy          `yaml:"policy"`
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
	c.Key = fmt.Sprintf("%s-%s-%s", environment.CI_MERGE_REQUEST_ID.Get(), path, pipelineType)
}
func NewCache() Cache {
	return Cache{
		Paths:     *wrapper.NewArray[string](),
		Policy:    cache.PullPush,
		Unprotect: true,
	}
}

func (c *Cache) GetRenderedValue() (*CacheYaml, error) {
	if c.Key == "" || !c.Active {
		return nil, nil
	}

	if c.Key == "" && len(c.Paths.Get()) > 0 {
		return nil, errors.New("no cache key but active paths found")
	}
	return &CacheYaml{
		Key:       c.Key,
		Paths:     c.Paths.Get(),
		Unprotect: c.Unprotect,
		Policy:    c.Policy.String(),
	}, nil
}

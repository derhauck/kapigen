package gitlab

import "kapigen.kateops.com/internal/pipeline/wrapper"

type CachePolicyEnum int

const (
	CachePolicyEnumPull CachePolicyEnum = iota
	CachePolicyEnumPush
	CachePolicyPullPush
)

func (c CachePolicyEnum) Policy() string {
	return []string{
		"pull",
		"push",
		"pull-push",
	}[c]
}

type CacheYaml struct {
	Key       string   `yaml:"key"`
	Paths     []string `yaml:"paths"`
	Unprotect bool     `yaml:"unprotect"`
	Policy    string   `yaml:"policy"`
}

type Cache struct {
	Key       string                `yaml:"key"`
	Paths     wrapper.Slice[string] `yaml:"paths"`
	Unprotect bool                  `yaml:"unprotect"`
	Policy    CachePolicyEnum       `yaml:"policy"`
}

func (c *Cache) getRenderedValue() CacheYaml {
	return CacheYaml{
		Key:       c.Key,
		Paths:     c.Paths.Get(),
		Unprotect: c.Unprotect,
		Policy:    c.Policy.Policy(),
	}
}

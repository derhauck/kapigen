package services

import (
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type Ci struct {
	name       docker.Image
	alias      string
	port       int32
	variables  map[string]string
	entrypoint wrapper.StringSlice
	command    wrapper.StringSlice
}

func (c *Ci) Entrypoint() *wrapper.StringSlice {
	return &c.entrypoint
}

func (c *Ci) Command() *wrapper.StringSlice {
	return &c.command
}

func (c *Ci) AddVariable(key string, value string) *Ci {
	c.variables[key] = value
	return c
}

func New(image docker.Image, alias string, port int32) *Ci {
	return &Ci{
		name:  image,
		alias: alias,
		port:  port,
	}
}

type Cis struct {
	Values *[]*Ci
}

func (c *Cis) Get() []*Ci {
	if c.Values == nil {
		c.Values = &[]*Ci{}
	}
	return *c.Values
}

func (c *Cis) Add(service *Ci) *Cis {
	values := append(c.Get(), service)
	c.Values = &values
	return c
}
func (c *Ci) Render() *Yaml {
	return &Yaml{
		Name:       c.name.Image(),
		Entrypoint: c.entrypoint.Get(),
		Command:    c.command.Get(),
		Alias:      c.alias,
		Variables:  c.variables,
	}
}

func (c *Cis) Render() *Yamls {
	if c != nil {
		var values = Yamls{}
		for _, ci := range c.Get() {
			values = append(values, ci.Render())
		}
		return &values
	}
	return nil
}

type Yaml struct {
	Name       string            `yaml:"name"`
	Entrypoint []string          `yaml:"entrypoint,omitempty"`
	Command    []string          `yaml:"command,omitempty"`
	Alias      string            `yaml:"alias"`
	Variables  map[string]string `yaml:"variables,omitempty"`
}

type Yamls []*Yaml

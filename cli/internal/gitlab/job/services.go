package job

import (
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type Service struct {
	name       string
	alias      string
	port       int32
	variables  map[string]string
	entrypoint wrapper.StringSlice
	command    wrapper.StringSlice
}

func (c *Service) Entrypoint() *wrapper.StringSlice {
	return &c.entrypoint
}

func (c *Service) Command() *wrapper.StringSlice {
	return &c.command
}

func (c *Service) AddVariable(key string, value string) *Service {
	if c.variables == nil {
		c.variables = map[string]string{}
	}
	c.variables[key] = value
	return c
}

func NewService(image string, alias string, port int32) *Service {
	return &Service{
		name:  image,
		alias: alias,
		port:  port,
	}
}

type Services struct {
	Values *[]*Service
}

func (c *Services) Get() []*Service {
	if c.Values == nil {
		c.Values = &[]*Service{}
	}
	return *c.Values
}

func (c *Services) Add(service *Service) *Services {
	values := append(c.Get(), service)
	c.Values = &values
	return c
}
func (c *Service) Render() *ServiceYaml {
	return &ServiceYaml{
		Name:       c.name,
		Entrypoint: c.entrypoint.Get(),
		Command:    c.command.Get(),
		Alias:      c.alias,
		Variables:  c.variables,
	}
}

func (c *Services) Render() *ServiceYamls {
	if c != nil {
		var values = ServiceYamls{}
		for _, ci := range c.Get() {
			values = append(values, ci.Render())
		}
		return &values
	}
	return nil
}

type ServiceYaml struct {
	Name       string            `yaml:"name"`
	Entrypoint []string          `yaml:"entrypoint,omitempty"`
	Command    []string          `yaml:"command,omitempty"`
	Alias      string            `yaml:"alias"`
	Variables  map[string]string `yaml:"variables,omitempty"`
}

type ServiceYamls []*ServiceYaml

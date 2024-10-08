package job

import (
	"fmt"

	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/images"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func (c *CiJob) SetStage(stage stages.Stage) *CiJob {
	c.Stage = stage

	return c
}

func (c *CiJob) SetCodeCoverage(regex string) *CiJob {
	c.Coverage = regex

	return c
}
func (c *CiJob) AddArtifact(artifact Artifacts) *CiJob {
	c.Artifact = artifact

	return c
}

func (c *CiJob) AddVariable(key string, value string) *CiJob {
	if c.Variables == nil {
		c.Variables = map[string]string{}
	}
	c.Variables[key] = value

	return c
}

func (c *CiJob) AddService(service *Service) {
	c.Services.Add(service)
}

func NewCiJob(imageName string) *CiJob {
	return &CiJob{
		Script: NewScript(),
		Cache:  NewCache(),
		Image: Image{
			Name:       imageName,
			PullPolicy: images.Always,
		},
		Stage:        stages.NewStage(),
		AfterScript:  NewAfterScript(),
		BeforeScript: NewBeforeScript(),
	}
}
func (c *CiJob) AddScript(script string) *CiJob {
	c.Script.Value.Push(script)

	return c
}
func (c *CiJob) AddScriptf(script string, a ...any) *CiJob {
	c.Script.Value.Push(fmt.Sprintf(script, a...))

	return c
}
func (c *CiJob) AddScripts(scripts []string) *CiJob {
	c.Script.Value.Push(scripts...)

	return c
}
func (c *CiJob) AddAfterScript(script string) *CiJob {
	c.AfterScript.Value.Push(script)

	return c
}
func (c *CiJob) AddBeforeScript(script string) *CiJob {
	c.BeforeScript.Value.Push(script)

	return c
}
func (c *CiJob) AddBeforeScriptf(script string, a ...any) *CiJob {
	c.BeforeScript.Value.Push(fmt.Sprintf(script, a...))

	return c
}
func (c *CiJob) AddBeforeScripts(scripts []string) *CiJob {
	c.BeforeScript.Value.Push(scripts...)

	return c
}

func (c *CiJob) SetImageName(image string) *CiJob {
	c.Image.Name = image

	return c
}

func (c *CiJob) SetImageEntrypoint(entrypoint wrapper.Array[string]) *CiJob {
	c.Image.Entrypoint = entrypoint

	return c
}

func (c *CiJob) SetImagePullPolicy(policy images.PullPolicy) *CiJob {
	c.Image.PullPolicy = policy

	return c
}

func (c *CiJob) TagMediumPressure() *CiJob {
	c.Tags.Add(enum.TagPressureMedium)

	return c
}

func (c *CiJob) TagHighPressure() *CiJob {
	c.Tags.Add(enum.TagPressureExclusive)

	return c
}

func (c *CiJob) AddSecret(name string, secret Secret) *CiJob {
	if c.Secrets == nil {
		c.Secrets = Secrets{}
	}
	if c.Secrets[name] != nil {
		logger.Infof("secret %s already exists, will be overwritten", name)
	}
	c.Secrets[name] = secret

	return c
}

func (c *CiJob) AddIdToken(name string, aud string) *CiJob {
	if c.IdTokens == nil {
		c.IdTokens = IdTokens{}
	}
	if c.IdTokens[name] != nil {
		logger.Infof("id token %s already exists, will be overwritten", name)
	}
	c.IdTokens[name] = &IdToken{
		AUD: aud,
	}

	return c
}

func (c *CiJob) AddEnvironment(environment *Environment) *CiJob {
	c.Environment = environment

	return c
}

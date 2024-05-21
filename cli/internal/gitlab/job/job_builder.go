package job

import (
	"kapigen.kateops.com/internal/gitlab/images"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func (c *CiJob) SetStage(stage stages.Stage) *CiJob {
	c.Stage = stage
	return c
}

func (c *CiJob) SetCodeCoverage(regex string) {
	c.Coverage = regex
}
func (c *CiJob) AddArtifact(artifact Artifact) *CiJob {
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
	c.Script.Value.Add(script)

	return c
}
func (c *CiJob) AddScriptf(script string, a ...any) *CiJob {
	c.Script.Value.Addf(script, a...)

	return c
}
func (c *CiJob) AddScripts(scripts []string) *CiJob {
	c.Script.Value.AddSlice(scripts)

	return c
}

func (c *CiJob) AddBeforeScript(script string) *CiJob {
	c.BeforeScript.Value.Add(script)

	return c
}
func (c *CiJob) AddBeforeScriptf(script string, a ...any) *CiJob {
	c.BeforeScript.Value.Addf(script, a...)

	return c
}
func (c *CiJob) AddBeforeScripts(scripts []string) *CiJob {
	c.BeforeScript.Value.AddSlice(scripts)

	return c
}

func (c *CiJob) SetImageName(image string) *CiJob {
	c.Image.Name = image

	return c
}

func (c *CiJob) SetImageEntrypoint(entrypoint wrapper.StringSlice) *CiJob {
	c.Image.Entrypoint = entrypoint

	return c
}

func (c *CiJob) SetImagePullPolicy(policy images.PullPolicy) *CiJob {
	c.Image.PullPolicy = policy

	return c
}

func (c *CiJob) TagMediumPressure() *CiJob {
	c.Tags.Add(tags.PRESSURE_MEDIUM)

	return c
}

func (c *CiJob) TagHighPressure() *CiJob {
	c.Tags.Add(tags.PRESSURE_BUILDKIT)

	return c
}

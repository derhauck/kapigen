package job

import (
	"kapigen.kateops.com/internal/gitlab/images"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func (j *CiJob) AddScript(script string) *CiJob {
	j.Script.Value.Add(script)

	return j
}

func (j *CiJob) AddScripts(scripts []string) *CiJob {
	j.Script.Value.AddSeveral(scripts)

	return j
}

func (j *CiJob) AddBeforeScript(script string) *CiJob {
	j.BeforeScript.Value.Add(script)

	return j
}

func (j *CiJob) AddBeforeScripts(scripts []string) *CiJob {
	j.BeforeScript.Value.AddSeveral(scripts)

	return j
}

func (j *CiJob) SetImageName(image string) *CiJob {
	j.Image.Name = image

	return j
}

func (j *CiJob) SetImageEntrypoint(entrypoint wrapper.StringSlice) *CiJob {
	j.Image.Entrypoint = entrypoint

	return j
}

func (j *CiJob) SetImagePullPolicy(policy images.PullPolicy) *CiJob {
	j.Image.PullPolicy = policy

	return j
}

func (j *CiJob) TagMediumPressure() *CiJob {
	j.Tags.Add(tags.PRESSURE_MEDIUM)

	return j
}

func (j *CiJob) TagHighPressure() *CiJob {
	j.Tags.Add(tags.PRESSURE_BUILDKIT)

	return j
}

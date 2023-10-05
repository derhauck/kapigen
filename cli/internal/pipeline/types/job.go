package types

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/logger"
)

type Job struct {
	Names       []string
	CiJob       *gitlab.CiJob
	Needs       Needs
	currentName int
	fn          []func(job *gitlab.CiJob)
	CiJobYaml   *gitlab.CiJobYaml
}

func (j *Job) AddNeed(job *Job) *Job {
	j.Needs = append(j.Needs, NewNeed(job))
	return j
}

func (j *Job) AddSeveralNeeds(needs *Needs) *Job {
	for _, need := range needs.GetNeeds() {
		if !j.HasNeed(need) {
			j.Needs = append(j.Needs, need)
		}
	}
	return j
}

func (j *Job) HasNeed(need *Need) bool {
	for _, availableNeed := range j.Needs {
		if availableNeed == need {
			return true
		}
	}
	return false
}

func (j *Job) RenderNeeds() *Job {
	if j.Needs != nil {
		j.CiJobYaml.Needs = j.Needs.NeedsYaml()
	}
	return j
}

func (j *Job) GetName() string {
	var name = j.Names[0]
	for i := 1; i < len(j.Names) && i < j.currentName; i++ {
		name = fmt.Sprintf("%s - %s", name, j.Names[i])
	}
	return name
}

func (j *Job) UniqueName() error {

	if j.currentName+1 <= len(j.Names) {
		j.currentName++
		logger.Info(fmt.Sprintf("added unique name for Job: %s", j.GetName()))
		return nil
	}
	return errors.New(fmt.Sprintf("job '%s' can not be more unique", j.GetName()))
}
func (j *Job) EvaluateName(jobs *Jobs) (*Job, error) {
	for _, job := range jobs.GetJobs() {
		if j != job {
			if j.compareConfiguration(job) {
				job.AddSeveralNeeds(job.EvaluateNeeds(&j.Needs))
				return nil, nil
			}
			if j.compare(job) {
				_, err := j.EvaluateName(jobs)
				if err != nil {
					logger.ErrorE(err)
					return nil, err
				}
			}
		}

	}

	return j, nil
}

func (j *Job) compare(job *Job) bool {
	if j.GetName() == job.GetName() {
		err := j.UniqueName()
		if err != nil {
			logger.ErrorE(err)
			return false
		}
		return true
	}

	return false
}

func (j *Job) compareConfiguration(job *Job) bool {
	if j.CiJobYaml.String() == job.CiJobYaml.String() {
		return true
	}
	return false
}

func (j *Job) AddName(name string) *Job {
	j.Names = append(j.Names, name)
	return j
}

func NewJob(name string, image docker.Image, fn func(job *gitlab.CiJob)) *Job {
	var newJob = &Job{
		Names:       []string{name},
		CiJob:       gitlab.NewCiJob(image),
		currentName: 0,
		fn:          []func(job *gitlab.CiJob){fn},
		Needs:       Needs{},
	}

	return newJob
}

func (j *Job) Render() *Job {
	for _, fn := range j.fn {
		fn(j.CiJob)
	}
	j.CiJobYaml = j.CiJob.Render(j.Needs.NeedsYaml())
	return j
}

type Jobs []*Job

func (j *Jobs) GetJobs() []*Job {
	return *j
}

func (j *Job) EvaluateNeeds(needs *Needs) *Needs {
	var finalNeeds Needs
	if needs != nil {
		for _, need := range needs.GetNeeds() {
			finalNeeds = append(finalNeeds, need)
		}
	}
	var jobNeeds = j.Needs
	if jobNeeds != nil {
		var tmp Needs
		tmp = append(jobNeeds.GetNeeds(), finalNeeds.GetNeeds()...)
		return &tmp
	}
	return nil
}

package types

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/logger"
)

type Job struct {
	Names       []string
	CiJob       *gitlab.CiJob
	Need        Jobs
	currentName int
	fn          []func(job *Job)
	CiJobYaml   *gitlab.CiJobYaml
}

func (j *Job) AddNeed(job *Job) {
	j.Need = append(j.Need, job)
}

func (j *Job) AddNeeds(job *Jobs) {
	j.Need = append(j.Need, job.GetJobs()...)
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
		logger.Info(fmt.Sprintf("added unique name for Job: %s", j.GetName()))
		j.currentName++
		return nil
	}
	return errors.New(fmt.Sprintf("job '%s' can not be more unique", j.GetName()))
}
func (j *Job) EvaluateName(jobs Jobs) (*Job, error) {
	for _, job := range jobs {
		if j.compareConfiguration(job) {
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

	return j, nil
}

func (j *Job) compare(job *Job) bool {
	if j == job {
		return false
	}
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

func NewJob(name string, fn func(job *Job)) *Job {
	var newJob = &Job{
		Names:       []string{name},
		CiJob:       gitlab.NewCiJob(),
		currentName: 0,
		fn:          []func(job *Job){fn},
		Need:        Jobs{},
	}

	return newJob
}

func (j *Job) Render() {
	for _, fn := range j.fn {
		fn(j)
	}
	j.CiJobYaml = j.CiJob.Render()
}

type Jobs []*Job

func (j *Jobs) GetJobs() []*Job {
	return *j
}

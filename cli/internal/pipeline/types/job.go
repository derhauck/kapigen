package types

import (
	"errors"
	"fmt"
	"github.com/kylelemons/godebug/diff"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
	"os"
)

type Job struct {
	Names       []string
	CiJob       *job.CiJob
	Needs       Needs
	currentName int
	fn          []func(job *job.CiJob)
	CiJobYaml   *job.CiJobYaml
}

func (j *Job) AddJobAsNeed(job *Job) *Job {
	j.Needs = append(j.Needs, NewNeed(job))
	return j
}

func (j *Job) AddNeed(need *Need) *Job {
	j.Needs.AddNeed(need)
	return j
}
func (j *Job) ReplaceNeed(old *Need, new *Need) *Job {
	j.Needs.ReplaceJob(old, new)
	return j
}

func (j *Job) RemoveNeed(need *Need) bool {
	return j.Needs.RemoveNeed(need)
}
func (j *Job) AddNeedByStage(job *Job, stage stages.Stage) *Job {
	if job.CiJob.Stage <= stage {
		j.AddJobAsNeed(job)
	}
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

func (j *Job) DynamicMerge(jobs *Jobs) (*Job, error) {
	for _, jobEvaluate := range jobs.GetJobs() {
		if j != jobEvaluate {
			if j.compareConfiguration(jobEvaluate) {
				jobEvaluate.AddSeveralNeeds(jobEvaluate.EvaluateNeeds(&j.Needs))
				for _, jobRemoveNeed := range jobs.GetJobs() {
					for _, need := range jobRemoveNeed.Needs.GetNeeds() {
						if need.Job == j {
							oldNeed := jobRemoveNeed.Needs.GetNeed(j)
							if oldNeed != nil {
								jobRemoveNeed.ReplaceNeed(oldNeed, NewNeed(jobEvaluate))
							}
						}
					}
				}
				return nil, nil
			}
		}
	}
	return j, nil
}

func (j *Job) EvaluateName(jobs *Jobs) (*Job, error) {
	for _, jobEvaluate := range jobs.GetJobs() {
		if j != jobEvaluate {
			if j.compare(jobEvaluate) {
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
	if os.Getenv("DEBUG_DIFF") != "" {
		logger.DebugAny(diff.Diff(j.CiJobYaml.String(), job.CiJobYaml.String()))
	}

	return false
}

func (j *Job) AddName(name string) *Job {
	j.Names = append(j.Names, name)
	return j
}

func NewJob(name string, image string, fn func(ciJob *job.CiJob)) *Job {
	var newJob = &Job{
		Names:       []string{name},
		CiJob:       job.NewCiJob(image),
		currentName: 0,
		fn:          []func(job *job.CiJob){fn},
		Needs:       Needs{},
	}

	return newJob
}

func (j *Job) Render() error {
	for _, fn := range j.fn {
		fn(j.CiJob)
	}
	var err error
	j.CiJobYaml, err = j.CiJob.Render(j.Needs.NeedsYaml())
	if err != nil {
		return fmt.Errorf("job='%s'  can not be rendered: %w", j.Names, err)
	}
	return nil
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

func (j *Jobs) DynamicMerge() (*Jobs, error) {
	var evaluatedJobs Jobs
	var jobsToEvaluate Jobs
	jobsToEvaluate = append(jobsToEvaluate, j.GetJobs()...)
	for _, currentJob := range j.GetJobs() {
		evaluatedJob, err := currentJob.DynamicMerge(&jobsToEvaluate)
		if err != nil {
			return nil, err
		}
		if evaluatedJob != nil {
			evaluatedJobs = append(evaluatedJobs, evaluatedJob)
		} else {
			var resizedJobsToEvaluate Jobs
			for i := range jobsToEvaluate {
				if jobsToEvaluate[i] == currentJob && i < len(jobsToEvaluate) {
					var tmp = jobsToEvaluate[i+1:]
					resizedJobsToEvaluate = append(jobsToEvaluate[:i], tmp...)
				}
			}
			jobsToEvaluate = resizedJobsToEvaluate
		}

	}
	return &evaluatedJobs, nil
}

func (j *Jobs) EvaluateNames() (*Jobs, error) {
	var evaluatedJobs Jobs
	var jobsToEvaluate Jobs
	jobsToEvaluate = append(jobsToEvaluate, j.GetJobs()...)
	for _, currentJob := range j.GetJobs() {
		evaluatedJob, err := currentJob.EvaluateName(&jobsToEvaluate)
		if err != nil {
			return nil, err
		}
		if evaluatedJob != nil {
			evaluatedJobs = append(evaluatedJobs, evaluatedJob)
		} else {
			var resizedJobsToEvaluate Jobs
			for i := range jobsToEvaluate {
				if jobsToEvaluate[i] == currentJob && i < len(jobsToEvaluate) {
					var tmp = jobsToEvaluate[i+1:]
					resizedJobsToEvaluate = append(jobsToEvaluate[:i], tmp...)
				}
			}
			jobsToEvaluate = resizedJobsToEvaluate
		}

	}
	return &evaluatedJobs, nil
}
func JobsToMap(jobs *Jobs) map[string]interface{} {
	var ciPipeline = make(map[string]interface{})
	for _, evaluatedJob := range *jobs {
		evaluatedJob.RenderNeeds()
		ciPipeline[evaluatedJob.GetName()] = evaluatedJob.CiJobYaml
	}
	return ciPipeline
}

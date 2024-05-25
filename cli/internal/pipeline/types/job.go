package types

import (
	"errors"
	"fmt"
	"os"

	"github.com/kylelemons/godebug/diff"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
)

type Job struct {
	Names        []string
	CiJob        *job.CiJob
	Needs        Needs
	currentName  int
	fn           []func(job *job.CiJob)
	CiJobYaml    *job.CiJobYaml
	PipelineId   string
	ExternalTags []string
}

func (j *Job) AddJobAsNeed(job *Job) *Job {
	j.Needs.AddNeed(NewNeed(job))
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
		if !j.Needs.HasNeed(need) {
			j.Needs.AddNeed(need)
		}
	}
	return j
}

func (j *Job) RenderNeeds() *Job {
	if j.CiJobYaml == nil {
		err := j.Render()
		if err != nil {
			return nil
		}
	}
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
				jobEvaluate.AddSeveralNeeds(&j.Needs)
				for _, jobRemoveNeed := range jobs.GetJobs() {
					jobRemoveNeed.Needs.ReplaceJob(j, NewNeed(jobEvaluate))
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

func (j *Job) toYamlConfiguration() *job.CiJobYaml {
	if j.CiJobYaml == nil {
		err := j.Render()
		if err != nil {
			logger.ErrorE(err)
			return nil
		}
	}
	return j.CiJobYaml
}

func (j *Job) compareConfiguration(job *Job) bool {
	if j.toYamlConfiguration().String() == job.toYamlConfiguration().String() {
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
		currentName: 2,
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
	j.CiJobYaml, err = j.CiJob.Render(j.Needs.NeedsYaml(), j.ExternalTags)
	if err != nil {
		return fmt.Errorf("job='%s'  can not be rendered: %w", j.Names, err)
	}
	return nil
}

type Jobs []*Job

func (j *Jobs) SetJobsAsNeed(jobs *Jobs) *Jobs {
	for _, currentJob := range j.GetJobs() {
		for _, currentNeed := range jobs.GetJobs() {
			currentJob.AddJobAsNeed(currentNeed)
		}
	}

	return j
}

func (j *Jobs) AddJob(job *Job) *Jobs {
	*j = append(*j, job)
	return j
}
func (j *Jobs) GetJobs() []*Job {
	return *j
}

func (j *Job) EvaluateNeeds(needs *Needs) {
	if needs != nil {
		for _, need := range needs.GetNeeds() {
			j.Needs.AddNeed(need)
		}
	}

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
func (j *Jobs) FindJobsByPipelineId(pipelineId string) (*Jobs, error) {
	found := Jobs{}
	for _, currentJob := range j.GetJobs() {
		if currentJob.PipelineId == pipelineId {
			found = append(found, currentJob)
		}
	}
	if len(found) == 0 {
		return &found, fmt.Errorf("can not find pipeline as need: %s", pipelineId)
	}
	return &found, nil
}
func JobsToMap(jobs *Jobs) map[string]interface{} {
	var ciPipeline = make(map[string]interface{})
	for _, evaluatedJob := range *jobs {
		evaluatedJob.RenderNeeds()
		ciPipeline[evaluatedJob.GetName()] = evaluatedJob.CiJobYaml
	}
	return ciPipeline
}

func (j *Jobs) OverwriteTags(tags []string) {
	if len(tags) > 0 && len(tags) > 0 {
		for _, evaluatedJob := range j.GetJobs() {
			if evaluatedJob.CiJobYaml != nil {
				evaluatedJob.CiJobYaml.Tags = tags
			}
		}
	}
}

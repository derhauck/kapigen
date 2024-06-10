package types

import (
	"kapigen.kateops.com/internal/gitlab/job"
)

type Need struct {
	Optional bool
	Job      *Job
}

func (n *Need) NotOptional() *Need {
	n.Optional = false
	return n
}

func (n *Need) HasJob(job *Job) bool {
	return n.Job == job
}

func (n *Needs) ReplaceJob(old *Job, new *Need) *Needs {
	if n.RemoveJob(old) {
		n.AddNeed(new)
	}
	return n
}

type Needs []*Need

func (n *Needs) GetNeeds() []*Need {
	return *n
}

func (n *Needs) GetNeed(job *Job) *Need {
	for _, need := range n.GetNeeds() {
		if need.HasJob(job) {
			return need
		}
	}
	return nil
}

func (n *Needs) HasNeed(need *Need) bool {
	if n.HasJob(need.Job) {
		return true
	}
	for _, availableNeed := range n.GetNeeds() {
		if availableNeed == need {
			return true
		}
	}

	return false
}

func (n *Needs) HasJob(job *Job) bool {
	for _, need := range n.GetNeeds() {
		if need.HasJob(job) {
			return true
		}
	}
	return false
}

func (n *Needs) AddNeed(need *Need) *Needs {
	if !n.HasJob(need.Job) {
		*n = append(*n, need)
	}
	return n
}

func (n *Needs) RemoveNeed(need *Need) bool {
	for i, availableNeed := range n.GetNeeds() {
		if availableNeed == need {
			tmp := n.GetNeeds()
			tmp = append(tmp[:i], tmp[i+1:]...)
			*n = tmp
			return true
		}
	}
	return false
}

func (n *Needs) RemoveJob(job *Job) bool {
	need := n.GetNeed(job)
	if need == nil {
		return false
	}
	for need != nil {
		n.RemoveNeed(need)
		need = n.GetNeed(job)
	}
	return true
}

func NewNeed(job *Job) *Need {
	return &Need{
		Optional: true,
		Job:      job,
	}
}

func (n *Need) Render() *job.NeedYaml {
	return &job.NeedYaml{
		Optional: n.Optional,
		Job:      n.Job.GetName(),
	}
}

func (n *Needs) NeedsYaml() *job.NeedsYaml {
	var needsYaml job.NeedsYaml
	for _, need := range n.GetNeeds() {
		needsYaml = append(needsYaml, need.Render())
	}
	return &needsYaml
}

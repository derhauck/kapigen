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
	if n.Job == job {
		return true
	}
	return false
}

func (n *Needs) ReplaceJob(old *Need, new *Need) *Needs {
	for _, need := range n.GetNeeds() {
		if need.HasJob(old.Job) {
			n.RemoveNeed(need)
			if !need.HasJob(new.Job) {
				n.AddNeed(new)
			}
		}
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

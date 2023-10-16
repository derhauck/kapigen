package types

import (
	"kapigen.kateops.com/internal/gitlab/job"
)

type Need struct {
	Optional bool
	Job      *Job
}

type Needs []*Need

func (n *Needs) GetNeeds() []*Need {
	return *n
}

func (n *Need) NotOptional() *Need {
	n.Optional = false
	return n
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

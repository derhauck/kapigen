package types

import "kapigen.kateops.com/internal/gitlab"

type Need struct {
	Optional bool
	Job      *Job
}

type Needs []*Need

func (n *Needs) GetNeeds() []*Need {
	return *n
}

func NewNeed(job *Job) *Need {
	return &Need{
		Optional: true,
		Job:      job,
	}
}

func (n *Need) Render() *gitlab.NeedYaml {
	return &gitlab.NeedYaml{
		Optional: n.Optional,
		Job:      n.Job.GetName(),
	}
}

func (n *Needs) NeedsYaml() *gitlab.NeedsYaml {
	var needsYaml gitlab.NeedsYaml
	for _, need := range n.GetNeeds() {
		needsYaml = append(needsYaml, need.Render())
	}
	return &needsYaml
}

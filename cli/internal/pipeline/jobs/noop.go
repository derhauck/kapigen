package jobs

import (
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/when"
)

func NewNoop() *types.Job {
	return types.NewJob("Noop", docker.Alpine_3_18.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(tags.PRESSURE_MEDIUM)
		ciJob.Script.Value.Push("echo \"successfully triggered\"")
		ciJob.Rules.Add(&job.Rule{
			When: job.NewWhen(when.OnSuccess),
		})
	})
}

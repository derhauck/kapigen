package jobs

import (
	"gitlab.com/kateops/kapigen/cli/internal/docker"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
)

func NewNoop() *types.Job {
	return types.NewJob("Noop", docker.Alpine_3_18.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(enum.TagPressureMedium)
		ciJob.Script.Value.Push("echo \"successfully triggered\"")
		ciJob.Rules.Add(&job.Rule{
			When: job.NewWhen(enum.WhenOnSuccess),
		})
	})
}

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
		ciJob.Script.Value.Add("echo \"successfully triggered\"")
		ciJob.Rules.Add(&job.Rule{
			If:   "($CI_MERGE_REQUEST_IID || $CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH)",
			When: job.NewWhen(when.OnSuccess),
		})
	})
}

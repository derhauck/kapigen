package jobs

import (
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/when"
)

func NewTag() *types.Job {
	return types.NewJob("Tag", docker.Alpine_3_18.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(tags.PRESSURE_MEDIUM)
		ciJob.Script.Value.Add("echo \"Hello World!\"")
		ciJob.Rules.Add(&job.Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: job.NewWhen(when.OnSuccess),
		})
		ciJob.Rules = *job.RulesNotKapigen(&ciJob.Rules)
	}).AddName("Default")
}
func NewTagKapigen() *types.Job {
	return types.NewJob("Tag", docker.GOLANG_1_21.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(tags.PRESSURE_MEDIUM)
		ciJob.BeforeScript.Value.Add("cd cli")
		ciJob.Script.Value.Add("go mod download").
			Add("go run . version -v --mode gitlab")
		ciJob.Rules.Add(&job.Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: job.NewWhen(when.OnSuccess),
		}).Add(&job.Rule{
			If:   "$CI_MERGE_REQUEST_ID",
			When: job.NewWhen(when.OnSuccess),
		})
		ciJob.Rules = *job.RulesKapigen(&ciJob.Rules)
	}).AddName("Kapigen")
}

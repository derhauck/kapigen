package jobs

import (
	"fmt"

	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/gitlab/tags"
	"kapigen.kateops.com/internal/logger/level"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
)

func NewTag(privateTokenName string) *types.Job {
	return types.NewJob("Versioning", docker.Kapigen_Latest.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(tags.PRESSURE_MEDIUM)
		ciJob.Stage = stages.FINAL
		ciJob.AddVariable("LOGGER_LEVEL", level.Info.String())
		ciJob.SetImageEntrypoint(*wrapper.NewStringSlice().Add(""))
		ciJob.Script.Value.Add(fmt.Sprintf("kapigen version new --mode gitlab --private-token '%s'", privateTokenName))
		ciJob.Rules.Add(&job.Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: job.NewWhen(when.OnSuccess),
		})
		ciJob.Rules = *job.RulesNotKapigen(&ciJob.Rules)
	}).AddName("Default")
}
func NewTagKapigen(privateTokenName string) *types.Job {
	return types.NewJob("Versioning", docker.GOLANG_1_21.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(tags.PRESSURE_MEDIUM)
		ciJob.Stage = stages.FINAL
		ciJob.AddVariable("LOGGER_LEVEL", level.Info.String())
		ciJob.BeforeScript.Value.Add("cd cli")
		ciJob.Script.Value.Add("go mod download").
			Add(fmt.Sprintf("go run . version new --mode gitlab --private-token '%s'", privateTokenName))
		ciJob.Rules.Add(&job.Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: job.NewWhen(when.OnSuccess),
		})
		ciJob.Rules = *job.RulesKapigen(&ciJob.Rules)
	}).AddName("Kapigen")
}

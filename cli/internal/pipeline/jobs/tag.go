package jobs

import (
	"fmt"

	"gitlab.com/kateops/kapigen/cli/internal/docker"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/logger/level"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func NewTag(privateTokenName string) *types.Job {
	return types.NewJob("Versioning", docker.Kapigen_Latest.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(enum.TagPressureMedium)
		ciJob.Stage = stages.FINAL
		ciJob.AddVariable("LOGGER_LEVEL", level.Info.String())
		ciJob.SetImageEntrypoint(*wrapper.NewArray[string]().Push(""))
		ciJob.Script.Value.Push(fmt.Sprintf("kapigen version new --mode gitlab --private-token '%s'", privateTokenName))
		ciJob.Rules.Add(&job.Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: job.NewWhen(enum.WhenOnSuccess),
		})
		ciJob.Rules = *job.RulesNotKapigen(&ciJob.Rules)
	}).AddName("Default")
}
func NewTagKapigen(privateTokenName string) *types.Job {
	return types.NewJob("Versioning", docker.GOLANG_1_21.String(), func(ciJob *job.CiJob) {
		ciJob.Tags.Add(enum.TagPressureMedium)
		ciJob.Stage = stages.FINAL
		ciJob.AddVariable("LOGGER_LEVEL", level.Info.String())
		ciJob.BeforeScript.Value.Push("cd cli")
		ciJob.Script.Value.Push("go mod download").
			Push(fmt.Sprintf("go run . version new --mode gitlab --private-token '%s'", privateTokenName))
		ciJob.Rules.Add(&job.Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: job.NewWhen(enum.WhenOnSuccess),
		})
		ciJob.Rules = *job.RulesKapigen(&ciJob.Rules)
	}).AddName("Kapigen")
}

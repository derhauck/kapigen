package golang

import (
	"fmt"

	"gitlab.com/kateops/kapigen/cli/internal/pipeline/types"
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job/artifact"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func Lint(imageName string, path string) *types.Job {
	return types.NewJob("Lint", imageName, func(ciJob *job.CiJob) {
		report := "junit.xml"
		reportPath := fmt.Sprintf("%s/%s", path, report)
		if path == "." {
			reportPath = report
		}
		ciJob.TagMediumPressure().
			SetStage(stages.TEST).
			AddBeforeScriptf("cd %s", path).
			AddScriptf("golangci-lint run -v --out-format=junit-xml:%s", report).
			AddArtifact(job.Artifacts{
				Name:  "report",
				Paths: *(wrapper.NewArray[string]().Push(reportPath)),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport(reportPath)),
				When: job.NewWhen(enum.WhenAlways),
			})
	})
}

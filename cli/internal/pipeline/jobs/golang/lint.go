package golang

import (
	"fmt"

	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/when"
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
			AddArtifact(job.Artifact{
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport(reportPath)),
				When: job.NewWhen(when.Always),
			})
	})
}

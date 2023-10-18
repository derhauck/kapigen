package golang

import (
	"fmt"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/job/artifact/reports"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func NewGolangTest(image string, path string) *types.Job {

	return types.NewJob("Unit Test", image, func(ciJob *job.CiJob) {
		var reportPath = fmt.Sprintf("%s/report.xml", path)
		if path == "." {
			reportPath = "report.xml"
		}
		ciJob.SetImageName(image).
			TagMediumPressure().
			AddBeforeScript(fmt.Sprintf("cd %s", path)).
			AddScript("go install github.com/jstemmer/go-junit-report/v2@latest").
			AddScript("go test -json ./... 2>&1 | go-junit-report -parser gojson -iocopy -out report.xml").
			AddVariable("KTC_PATH", path).
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *(wrapper.NewStringSlice().Add(reportPath)),
				Reports: artifact.Reports{
					CoverageReport: artifact.CoverageReport{
						Path:           reportPath,
						CoverageFormat: reports.Cobertura,
					},
					Junit: artifact.NewJunitReport(reportPath),
				},
			})

		ciJob.Rules = *job.DefaultPipelineRules()
	})
}

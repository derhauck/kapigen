package golang

import (
	"errors"
	"fmt"
	"strings"

	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/job/artifact/reports"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func NewUnitTest(imageName string, path string, packages []string, source string) (*types.Job, error) {
	if imageName == "" {
		return nil, errors.New("no imageName set, required")
	}
	if path == "" {
		return nil, errors.New("no path set, required")
	}

	return types.NewJob("Unit Test", imageName, func(ciJob *job.CiJob) {
		var reportPath = fmt.Sprintf("%s/report.xml", path)
		if path == "." {
			reportPath = "report.xml"
		}
		coveragePkg := strings.Join(packages, ",")
		testCmd := fmt.Sprintf("go test -json -cover %s -coverpkg=%s -coverprofile=profile.cov", source, coveragePkg)
		ciJob.TagMediumPressure().
			SetStage(stages.TEST).
			AddBeforeScriptf("cd %s", path).
			AddScript("go install github.com/jstemmer/go-junit-report/v2@latest").
			AddScriptf("%s 2>&1 | go-junit-report -parser gojson -iocopy -out report.xml || (go tool cover -func profile.cov; exit 1)", testCmd).
			AddScript("go tool cover -func profile.cov").
			AddVariable("KTC_PATH", path).
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *(wrapper.NewStringSlice().Add(reportPath)),
				Reports: artifact.NewReports().
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, reportPath)).
					SetJunit(artifact.NewJunitReport(reportPath)),
			}).
			SetCodeCoverage(`/\(statements\)(?:\s+)?(\d+(?:\.\d+)?%)/`)

	}), nil
}

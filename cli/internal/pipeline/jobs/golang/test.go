package golang

import (
	"errors"
	"fmt"
	"strings"

	"gitlab.com/kateops/kapigen/cli/internal/pipeline/types"
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job/artifact"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job/artifact/reports"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func NewUnitTest(imageName string, path string, packages []string, source string) (*types.Job, error) {
	if imageName == "" {
		return nil, errors.New("no imageName set, required")
	}
	if path == "" {
		return nil, errors.New("no path set, required")
	}

	return types.NewJob("Unit Test", imageName, func(ciJob *job.CiJob) {
		var reportPath = fmt.Sprintf("%s/junit.xml", path)
		if path == "." {
			reportPath = "report.xml"
		}
		coveragePkg := strings.Join(packages, ",")
		testCmd := fmt.Sprintf("go test -json -cover %s -coverpkg=%s -coverprofile=profile.cov", source, coveragePkg)
		ciJob.TagMediumPressure().
			SetStage(stages.TEST).
			AddBeforeScriptf("cd %s", path).
			AddScript("go install github.com/jstemmer/go-junit-report/v2@latest").
			AddScriptf("%s 2>&1 | go-junit-report -parser gojson -iocopy -out junit.xml || (go tool cover -func profile.cov; exit 1)", testCmd).
			AddScript("go tool cover -func profile.cov").
			AddVariable("KTC_PATH", path).
			AddArtifact(job.Artifacts{
				Name:  "report",
				Paths: *(wrapper.NewArray[string]().Push(reportPath)),
				Reports: artifact.NewReports().
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, reportPath)).
					SetJunit(artifact.NewJunitReport(reportPath)),
				When: job.NewWhen(enum.WhenAlways),
			}).
			SetCodeCoverage(`/\(statements\)(?:\s+)?(\d+(?:\.\d+)?%)/`)

	}), nil
}

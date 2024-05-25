package php

import (
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/job/artifact/reports"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func NewPhpUnit(imageName string, composerPath string, composerArgs string, phpUnitXmlPath string, phpUnitArgs string) (*types.Job, error) {

	return types.NewJob("Unit Test", imageName, func(ciJob *job.CiJob) {
		//var reportPath = fmt.Sprintf("%s/report.xml", composerPath)
		//if composerPath == "." {
		//	reportPath = "report.xml"
		//}
		ciJob.
			TagMediumPressure().
			SetStage(stages.TEST).
			AddVariable("XDEBUG_MODE", "coverage").
			AddScriptf("composer install --no-progress --working-dir=%s %s", composerPath, composerArgs).
			AddScriptf("php %s/vendor/bin/phpunit -c %s/phpunit.xml --log-junit report.xml  --coverage-text --colors=never --coverage-cobertura=coverage.cobertura.xml %s", composerPath, phpUnitXmlPath, phpUnitArgs).
			SetCodeCoverage(`/^\s*Lines:\s*\d+.\d+\%/`).
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *wrapper.NewStringSlice().Add("report.xml"),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport("report.xml")).
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, "coverage.cobertura.xml")),
			})
	}), nil
}

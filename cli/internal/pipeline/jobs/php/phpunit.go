package php

import (
	"fmt"

	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func NewPhpUnit(imageName string, path string) (*types.Job, error) {

	return types.NewJob("PHP Unit", imageName, func(ciJob *job.CiJob) {
		var reportPath = fmt.Sprintf("%s/report.xml", path)
		if path == "." {
			reportPath = "report.xml"
		}

		ciJob.TagMediumPressure().
			SetStage(stages.TEST).
			AddBeforeScriptf("cd %s", path).
			AddScriptf("composer install").
			AddScript("vendor/bin/phpunit --log-junit report.xml").
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *wrapper.NewStringSlice().Add(reportPath),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport(reportPath)),
			})
	}), nil
}

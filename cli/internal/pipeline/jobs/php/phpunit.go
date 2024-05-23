package php

import (
	"fmt"

	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func NewPhpUnit(imageName string, composerPath string, composerArgs string, phpUnitXmlPath string, phpUnitArgs string) (*types.Job, error) {

	return types.NewJob("Unit Test", imageName, func(ciJob *job.CiJob) {
		var reportPath = fmt.Sprintf("%s/report.xml", composerPath)
		if composerPath == "." {
			reportPath = "report.xml"
		}

		ciJob.TagMediumPressure().
			SetStage(stages.TEST).
			AddBeforeScriptf("cd %s", composerPath).
			AddScriptf("composer install --no-progress %s", composerArgs).
			AddScriptf("php vendor/bin/phpunit -c %s/phpunit.xml --log-junit report.xml %s", phpUnitXmlPath, phpUnitArgs).
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *wrapper.NewStringSlice().Add(reportPath),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport(reportPath)),
			})
	}), nil
}

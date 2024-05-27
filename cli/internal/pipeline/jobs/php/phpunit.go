package php

import (
	"fmt"

	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/job/artifact"
	"kapigen.kateops.com/internal/gitlab/job/artifact/reports"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

func NewPhpUnit(imageName string, composerPath string, composerArgs string, phpUnitXmlPath string, phpUnitArgs string, listenerPorts map[string]int32) (*types.Job, error) {

	return types.NewJob("Unit Test", imageName, func(ciJob *job.CiJob) {
		ciJob.
			TagMediumPressure().
			SetStage(stages.TEST).
			AddVariable("XDEBUG_MODE", "coverage").
			AddBeforeScriptf("composer install --working-dir=%s %s", composerPath, composerArgs).
			AddScript("while [ ! -f ${CI_PROJECT_DIR}/.ready ]; do sleep 1; done;").
			AddScriptf("php %s/vendor/bin/phpunit -c %s/phpunit.xml --log-junit report.xml  --coverage-text --colors=never --coverage-cobertura=coverage.cobertura.xml %s", composerPath, phpUnitXmlPath, phpUnitArgs).
			SetCodeCoverage(`/^\s*Lines:\s*\d+.\d+\%/`).
			AddAfterScript("tail ${CI_PROJECT_DIR}/.status").
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *wrapper.NewStringSlice().Add("report.xml"),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport("report.xml")).
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, "coverage.cobertura.xml")),
			})

		listener := job.NewService(docker.Alpine_3_18.String(), "kapigen-listener", 0)
		listener.Entrypoint().AddSeveral("sh", "-c")
		command := ""
		for name, port := range listenerPorts {
			command += fmt.Sprintf("while ! nc -vz -w 1 %s %d &> /dev/null; do echo \"waiting for %s\" >> ${CI_PROJECT_DIR}/.status; sleep 1; done; ", name, port, name)
		}
		command += "echo \"done\" > ${CI_PROJECT_DIR}/.ready"
		listener.Command().Add(command)
		ciJob.AddService(listener)
	}), nil
}

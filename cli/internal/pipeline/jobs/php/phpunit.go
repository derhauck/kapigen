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
	types2 "kapigen.kateops.com/internal/types"
)

func NewPhpUnit(imageName string, composerPath string, composerArgs string, phpunitXmlPath string, phpunitArgs string, phpUnitBin string, listenerPorts map[string]int32) (*types.Job, error) {

	if imageName == "" {
		return nil, types2.NewMissingArgError("imageName")
	}

	if composerPath == "" {
		return nil, types2.NewMissingArgError("composerPath")
	}
	return types.NewJob("Unit Test", imageName, func(ciJob *job.CiJob) {
		ciJob.
			TagMediumPressure().
			SetStage(stages.TEST).
			AddVariable("XDEBUG_MODE", "coverage").
			AddBeforeScriptf("composer install --working-dir=%s %s", composerPath, composerArgs).
			AddScript("while [ ! -f ${CI_PROJECT_DIR}/.ready ]; do sleep 1; done;").
			AddScriptf("php %s -c %s/phpunit.xml --log-junit junit.xml  --coverage-text --colors=never --coverage-cobertura=coverage.cobertura.xml %s", phpUnitBin, phpunitXmlPath, phpunitArgs).
			SetCodeCoverage(`/^\s*Lines:\s*\d+.\d+\%/`).
			AddAfterScript("tail ${CI_PROJECT_DIR}/.status").
			AddArtifact(job.Artifact{
				Name:  "report",
				Paths: *wrapper.NewArray[string]().Push("junit.xml"),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport("junit.xml")).
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, "coverage.cobertura.xml")),
			})

		listener := job.NewService(docker.Alpine_3_18.String(), "kapigen-listener", 0)
		listener.Entrypoint().Push("sh", "-c")
		command := ""
		for name, port := range listenerPorts {
			command += fmt.Sprintf("while ! nc -vz -w 1 %s %d &> /dev/null; do echo \"waiting for %s\" >> ${CI_PROJECT_DIR}/.status; sleep 1; done; ", name, port, name)
		}
		command += "while :; do echo \"done\" > ${CI_PROJECT_DIR}/.ready; sleep 10; done"
		listener.Command().Push(command)
		ciJob.AddService(listener)
	}), nil
}

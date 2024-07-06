package php

import (
	"fmt"

	"gitlab.com/kateops/kapigen/cli/internal/docker"
	"gitlab.com/kateops/kapigen/cli/types"
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job/artifact"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job/artifact/reports"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

func NewPhpUnit(imageName string, composerPath string, composerArgs string, phpunitXmlPath string, phpunitArgs string, phpUnitBin string, listenerPorts map[string]int32) (*types.Job, error) {

	if imageName == "" {
		return nil, wrapper.NewMissingArgError("imageName")
	}

	if composerPath == "" {
		return nil, wrapper.NewMissingArgError("composerPath")
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
			AddArtifact(job.Artifacts{
				Name:  "report",
				Paths: *wrapper.NewArray[string]().Push("junit.xml"),
				Reports: artifact.NewReports().
					SetJunit(artifact.NewJunitReport("junit.xml")).
					SetCoverageReport(artifact.NewCoverageReport(reports.Cobertura, "coverage.cobertura.xml")),
				When: job.NewWhen(enum.WhenAlways),
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

package golang

import (
	"fmt"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewGolangTest(image string, path string) *types.Job {

	return types.NewJob("Unit Test", image, func(ciJob *job.CiJob) {
		ciJob.SetImageName(image).
			TagMediumPressure().
			AddBeforeScript(fmt.Sprintf("cd %s", path)).
			AddScript("go install github.com/jstemmer/go-junit-report/v2@latest").
			AddScript("go test -json ./... 2>&1 | go-junit-report -parser gojson -iocopy -out report.xml").
			AddVariable("KTC_PATH", path)
		//AddArtifact(job.Artifact{
		//	Name:  "report",
		//	Paths: *(wrapper.NewStringSlice().Add("TEST")),
		//})

		ciJob.Rules = *job.DefaultPipelineRules()
	})
}

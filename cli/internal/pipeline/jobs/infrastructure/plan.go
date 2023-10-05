package infrastructure

import (
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformPlan(path string, state string, s3 bool) *types.Job {
	return types.NewJob("Plan", docker.Terraform_Base, func(job *gitlab.CiJob) {
		job.Script.Value.AddSeveral([]string{
			"echo \"test\"",
			//fmt.Sprintf(`echo "state: %s"`, state),
			//fmt.Sprintf(`echo "s3: %v"`, s3),
		})
		job.AllowFailure.ExitCodes.Add(300)

		job.Rules = *rules.DefaultReleasePipelineRules()

	})
}

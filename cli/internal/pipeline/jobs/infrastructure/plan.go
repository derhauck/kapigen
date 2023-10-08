package infrastructure

import (
	"fmt"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/gitlab/environment"
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformPlan(path string, state string, s3 bool) *types.Job {
	return types.NewJob("Plan", docker.Terraform_Base, func(job *gitlab.CiJob) {
		job.Script.Value.
			Add(fmt.Sprintf("echo \"%s\"", state)).
			Add("terraform plan")

		project, err := environment.Lookup(environment.CI_PROJECT_ID)
		if err != nil {
			logger.ErrorE(err)
			project = "test"
		}
		job.Variables = map[string]string{}
		job.Stage = stages.BUILD
		job.Variables["TF_STATE_PROJECT"] = project
		job.Rules = *rules.DefaultReleasePipelineRules()
		job.Cache.Paths.Add(fmt.Sprintf("%s/.terraform", path))
		job.Cache.SetActive().
			SetPolicy(gitlab.CachePolicyEnumPull)

	})
}

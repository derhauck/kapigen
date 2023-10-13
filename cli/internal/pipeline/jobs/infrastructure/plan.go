package infrastructure

import (
	"fmt"
	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/gitlab/cache"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformPlan(path string, state string, s3 bool) *types.Job {
	return types.NewJob("Plan", docker.Terraform_Base, func(ciJob *gitlab.CiJob) {
		ciJob.Script.Value.
			Add(fmt.Sprintf("echo \"%s\"", state)).
			Add("terraform plan")

		project, err := environment.CI_PROJECT_ID.Lookup()
		if err != nil {
			logger.ErrorE(err)
			project = "test"
		}
		ciJob.Variables = map[string]string{}
		ciJob.Stage = stages.BUILD
		ciJob.Variables["TF_STATE_PROJECT"] = project
		ciJob.Rules = *job.DefaultReleasePipelineRules()
		ciJob.Cache.Paths.Add(fmt.Sprintf("%s/.terraform", path))
		ciJob.Cache.SetActive().
			SetPolicy(cache.Pull)

	})
}

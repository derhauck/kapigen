package infrastructure

import (
	"fmt"

	"kapigen.kateops.com/internal/docker"
	"kapigen.kateops.com/internal/environment"
	"kapigen.kateops.com/internal/gitlab/cache"
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformPlan(path string, state string, s3 bool) *types.Job {
	return types.NewJob("Plan", docker.Terraform_Base.String(), func(ciJob *job.CiJob) {
		ciJob.Script.Value.
			Push(fmt.Sprintf("echo \"%s\"", state)).
			Push("terraform plan")

		project, err := environment.CI_PROJECT_ID.Lookup()
		if err != nil {
			logger.ErrorE(err)
			project = "test"
		}
		ciJob.Variables = map[string]string{}
		ciJob.Stage = stages.BUILD
		ciJob.Variables["TF_STATE_PROJECT"] = project
		ciJob.Rules = *job.DefaultOnlyReleasePipelineRules()
		ciJob.Cache.Paths.Push(fmt.Sprintf("%s/.terraform", path))
		ciJob.Cache.SetActive().
			SetPolicy(cache.Pull)

	})
}

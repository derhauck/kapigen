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

func NewTerraformInit(path string, state string, s3 bool) *types.Job {
	return types.NewJob("Init", docker.Terraform_Base.Image(), func(ciJob *job.CiJob) {
		ciJob.BeforeScript.Value.AddSeveral([]string{
			"echo \"credentials \\\\\"${CI_SERVER_HOST}\\\\\" {\\n  token = \\\\\"${CI_PIPELINE_TOKEN}\\\\\"\\n}\" > gitlab.tfrc",
			"export TF_CLI_CONFIG_FILE=${PWD}/gitlab.tfrc",
		})
		project, err := environment.CI_PROJECT_ID.Lookup()
		if err != nil {
			logger.ErrorE(err)
			project = "test"
		}

		ciJob.Stage = stages.INIT
		ciJob.Script.Value.Add(
			"terraform init \\\n" +
				" -backend-config=\"region=eu-central-1\" \\\n" +
				" -backend-config=\"access_key=${TF_STATE_ACCESS_KEY}\" \\\n" +
				" -backend-config=\"secret_key=${TF_STATE_SECRET_KEY}\" \\\n" +
				" -backend-config=\"bucket=${TF_STATE_BUCKET}\" \\\n" +
				fmt.Sprintf(` -backend-config="key=${TF_STATE_BUCKET}/states/%s/%s/terraform.tfstate"`, project, state),
		)

		ciJob.Variables = map[string]string{}
		ciJob.Variables["TF_STATE_PROJECT"] = project
		//job.AllowFailure.Failure = true
		ciJob.Rules = *job.DefaultPipelineRules()
		ciJob.Cache.SetPolicy(cache.PullPush).
			SetActive()
	})
}

package infrastructure

import (
	"fmt"
	"kapigen.kateops.com/docker"
	"kapigen.kateops.com/internal/gitlab"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformInit(state string, s3 bool) *types.Job {
	return types.NewJob("Init", docker.Terraform_Base, func(job *gitlab.CiJob) {
		job.Script.Value.AddSeveral([]string{
			"echo \"test\"",
			fmt.Sprintf(`echo "state: %s"`, state),
			fmt.Sprintf(`echo "s3: %v"`, s3),
		})
		job.Variables = map[string]string{}
		job.Variables["TEST"] = "yolo"
		//job.AllowFailure.Failure = true

	})
}

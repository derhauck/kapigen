package infrastructure

import (
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformInit(state string, s3 bool) *types.Job {
	return types.NewJob("Init", func(job *types.Job) {
		job.CiJob.Script.Value.AddSeveral([]string{
			"echo \"test\"",
		})
	})
}

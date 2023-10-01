package infrastructure

import (
	"fmt"
	"kapigen.kateops.com/internal/pipeline/types"
)

func NewTerraformPlan(state string, s3 bool) *types.Job {
	return types.NewJob("Plan", func(job *types.Job) {
		job.CiJob.Script.Value.AddSeveral([]string{
			"echo \"test\"",
			fmt.Sprintf(`echo "state: %s"`, state),
			fmt.Sprintf(`echo "s3: %v"`, s3),
		})
	})
}

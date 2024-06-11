package generic

import (
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/types"
	types2 "kapigen.kateops.com/internal/types"
)

func NewGenericJob(imageName string, stage stages.Stage, scripts []string) (*types.Job, error) {
	if imageName == "" {
		return nil, types2.NewMissingArgError("imageName")
	}

	return types.NewJob("Generic Job", imageName, func(ciJob *job.CiJob) {
		ciJob.SetStage(stage).
			AddScripts(scripts)
	}), nil
}

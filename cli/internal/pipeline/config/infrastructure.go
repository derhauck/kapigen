package config

import (
	"kapigen.kateops.com/factory"
	"kapigen.kateops.com/internal/pipeline/jobs/infrastructure"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Infrastructure struct {
	State string `yaml:"state"`
	S3    bool   `yaml:"s3"`
	Path  string `yaml:"path"`
}

func (i *Infrastructure) New() types.PipelineConfigInterface {
	return &Infrastructure{}
}

func (i *Infrastructure) Validate() error {
	if i.State == "" {
		i.State = "set-by-validation"
	}
	return nil
}

func (i *Infrastructure) Build(_ *factory.MainFactory, pipelineType types.PipelineType, _ string) (*types.Jobs, error) {
	var init = infrastructure.
		NewTerraformInit(i.Path, i.State, i.S3)
	var plan = infrastructure.
		NewTerraformPlan(i.Path, i.State, i.S3).
		AddJobAsNeed(init)
	var tmp = types.Jobs{init, plan}
	for _, job := range tmp {
		job.CiJob.Cache.SetDefaultCacheKey(i.Path, string(pipelineType))
	}
	return &tmp, nil
}

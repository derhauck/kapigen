package config

import (
	"kapigen.kateops.com/internal/pipeline/jobs/infrastructure"
	"kapigen.kateops.com/internal/pipeline/types"
)

type Infrastructure struct {
	State string `yaml:"state"`
	S3    bool   `yaml:"s3"`
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

func (i *Infrastructure) Build(path string, pipelineType types.PipelineType, Id string) (*types.Jobs, error) {
	var init = infrastructure.
		NewTerraformInit(path, i.State, i.S3)
	var plan = infrastructure.
		NewTerraformPlan(path, i.State, i.S3).
		AddNeed(init)
	var tmp = types.Jobs{init, plan}
	return &tmp, nil
}

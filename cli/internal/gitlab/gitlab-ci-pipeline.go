package gitlab

import (
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/pipeline/wrapper"
)

type CiPipelineDefault struct {
	AfterScript  AfterScript
	BeforeScript BeforeScript
}

type CiPipelineWorkflow struct {
	Name  string       `yaml:"name"`
	Rules *rules.Rules `yaml:"rules"`
}

type CiPipeline struct {
	Stages       wrapper.Slice[string] `json:"stages"`
	Workflow     CiPipelineWorkflow    `json:"workflow"`
	AllowFailure AllowFailure
	Default      CiPipelineDefault
}

type CiPipelineWorkflowYaml struct {
	Name  string           `yaml:"name"`
	Rules *rules.RulesYaml `yaml:"rules"`
}
type CiPipelineDefaultYaml struct {
	AfterScript  []string `yaml:"after_script"`
	BeforeScript []string `yaml:"before_script"`
}
type CiPipelineJson struct {
	Default  CiPipelineDefaultYaml  `yaml:"default"`
	Workflow CiPipelineWorkflowYaml `yaml:"workflow"`
}

func NewCiPipelineJson(pipeline *CiPipeline) *CiPipelineJson {
	return &CiPipelineJson{
		Default: CiPipelineDefaultYaml{
			AfterScript:  pipeline.Default.AfterScript.Value.Get(),
			BeforeScript: pipeline.Default.BeforeScript.Value.Get(),
		},
		Workflow: CiPipelineWorkflowYaml{
			Name:  pipeline.Workflow.Name,
			Rules: pipeline.Workflow.Rules.GetRenderedValue(),
		},
	}
}

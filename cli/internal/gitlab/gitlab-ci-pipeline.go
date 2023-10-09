package gitlab

import (
	"kapigen.kateops.com/internal/gitlab/rules"
	"kapigen.kateops.com/internal/gitlab/stages"
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
	Stages       *wrapper.StringSlice `yaml:"stages"`
	Workflow     CiPipelineWorkflow   `yaml:"workflow,omitempty"`
	AllowFailure AllowFailure         `yaml:"allow_failure,omitempty"`
	Default      CiPipelineDefault    `yaml:"default,omitempty"`
	Variables    map[string]string    `yaml:"variables,omitempty"`
}

func (c *CiPipeline) Render() *CiPipelineYaml {
	return NewCiPipelineYaml(c)
}
func NewDefaultCiPipeline() *CiPipeline {
	return &CiPipeline{
		Stages: wrapper.NewStringSlice().AddSeveral(stages.GetStages()),
		Default: CiPipelineDefault{
			AfterScript:  AfterScript{},
			BeforeScript: BeforeScript{},
		},
		Workflow: CiPipelineWorkflow{
			Name: "default",
			Rules: &rules.Rules{
				&rules.Rule{
					If:   "$CI_MERGE_REQUEST_ID",
					When: rules.NewWhen(rules.WhenEnumTypeAlways),
				},
				&rules.Rule{
					If:   "$CI_MERGE_REQUEST_IID && $CI_PIPELINE_SOURCE == 'merge_request_event'",
					When: rules.NewWhen(rules.WhenEnumTypeAlways),
				},
			},
		},
		Variables: map[string]string{
			"KTC_STOP_PIPELINE": "false",
		},
	}
}

type CiPipelineWorkflowYaml struct {
	Name  string                  `yaml:"name"`
	Rules *rules.RuleWorkflowYaml `yaml:"rules"`
}
type CiPipelineDefaultYaml struct {
	AfterScript  []string `yaml:"after_script"`
	BeforeScript []string `yaml:"before_script"`
}
type CiPipelineYaml struct {
	Default   CiPipelineDefaultYaml  `yaml:"default"`
	Workflow  CiPipelineWorkflowYaml `yaml:"workflow"`
	Stages    []string               `yaml:"stages"`
	Variables map[string]string      `yaml:"variables"`
}

func NewCiPipelineYaml(pipeline *CiPipeline) *CiPipelineYaml {
	return &CiPipelineYaml{
		Default: CiPipelineDefaultYaml{
			AfterScript:  pipeline.Default.AfterScript.Value.Get(),
			BeforeScript: pipeline.Default.BeforeScript.Value.Get(),
		},
		Workflow: CiPipelineWorkflowYaml{
			Name:  pipeline.Workflow.Name,
			Rules: pipeline.Workflow.Rules.GetRenderedWorkflowValue(),
		},
		Stages:    pipeline.Stages.Get(),
		Variables: pipeline.Variables,
	}
}

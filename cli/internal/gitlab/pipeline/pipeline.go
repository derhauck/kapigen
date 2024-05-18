package pipeline

import (
	"kapigen.kateops.com/internal/gitlab/job"
	"kapigen.kateops.com/internal/gitlab/stages"
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
)

type CiPipelineDefault struct {
	AfterScript  job.AfterScript
	BeforeScript job.BeforeScript
}

type CiPipelineWorkflow struct {
	Name  string     `yaml:"name"`
	Rules *job.Rules `yaml:"rules"`
}

type CiPipeline struct {
	Stages       *wrapper.StringSlice `yaml:"stages"`
	Workflow     CiPipelineWorkflow   `yaml:"workflow,omitempty"`
	AllowFailure job.AllowFailure     `yaml:"allow_failure,omitempty"`
	Default      CiPipelineDefault    `yaml:"default,omitempty"`
	Variables    map[string]string    `yaml:"variables,omitempty"`
}

func (c *CiPipeline) Render() *CiPipelineYaml {
	return NewCiPipelineYaml(c)
}

func NewDefaultCiPipeline() *CiPipeline {
	return &CiPipeline{
		Stages: wrapper.NewStringSlice().AddSlice(stages.GetAllStages()),
		Default: CiPipelineDefault{
			AfterScript:  job.NewAfterScript(),
			BeforeScript: job.NewBeforeScript(),
		},
		Workflow: CiPipelineWorkflow{
			Name: "default",
			Rules: &job.Rules{
				&job.Rule{
					If:   "$CI_MERGE_REQUEST_ID",
					When: job.NewWhen(when.Always),
				},
				&job.Rule{
					If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
					When: job.NewWhen(when.Always),
				},
				&job.Rule{
					If:   "$CI_MERGE_REQUEST_IID && $CI_PIPELINE_SOURCE == 'merge_request_event'",
					When: job.NewWhen(when.Always),
				},
				&job.Rule{
					If:   "$CI_COMMIT_TAG != null",
					When: job.NewWhen(when.Always),
				},
			},
		},
		Variables: map[string]string{
			"KTC_STOP_PIPELINE": "false",
		},
	}
}

type CiPipelineWorkflowYaml struct {
	Name  string                `yaml:"name"`
	Rules *job.RuleWorkflowYaml `yaml:"rules"`
}
type CiPipelineDefaultYaml struct {
	AfterScript  []string `yaml:"after_script,omitempty"`
	BeforeScript []string `yaml:"before_script,omitempty"`
}
type CiPipelineYaml struct {
	Default   CiPipelineDefaultYaml  `yaml:"default,omitempty"`
	Workflow  CiPipelineWorkflowYaml `yaml:"workflow,omitempty"`
	Stages    []string               `yaml:"stages"`
	Variables map[string]string      `yaml:"variables,omitempty"`
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
func (c *CiPipelineYaml) AddToMap(parentMap map[string]interface{}) {
	parentMap["default"] = c.Default
	parentMap["workflow"] = c.Workflow
	parentMap["stages"] = c.Stages
	parentMap["variables"] = c.Variables
}

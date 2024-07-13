package pipeline

import (
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type CiPipelineDefault struct {
	AfterScript  job.AfterScript
	BeforeScript job.BeforeScript
}

func (c *CiPipelineDefault) Validate() error {

	return nil
}

type CiPipelineWorkflow struct {
	Name  string     `yaml:"name"`
	Rules *job.Rules `yaml:"rules,omitempty"`
}

func (c *CiPipelineWorkflow) Validate() error {
	if c.Rules == nil {
		c.Rules = &job.Rules{}
	}

	return nil
}

type CiPipeline struct {
	Stages       *wrapper.Array[string] `yaml:"stages"`
	Workflow     *CiPipelineWorkflow    `yaml:"workflow,omitempty"`
	AllowFailure job.AllowFailure       `yaml:"allow_failure,omitempty"`
	Default      *CiPipelineDefault     `yaml:"default,omitempty"`
	Variables    map[string]string      `yaml:"variables,omitempty"`
}

func (c *CiPipeline) Validate() error {
	if c.Stages == nil {
		c.Stages = wrapper.NewArray[string]()
	}

	if c.Workflow == nil {
		c.Workflow = &CiPipelineWorkflow{}
	}

	if c.Default == nil {
		c.Default = &CiPipelineDefault{
			AfterScript:  job.NewAfterScript(),
			BeforeScript: job.NewBeforeScript(),
		}
	}

	err := c.Workflow.Validate()
	if err != nil {
		return err
	}

	err = c.Default.Validate()
	if err != nil {
		return err
	}
	return nil
}

func (c *CiPipeline) Render() (*CiPipelineYaml, error) {
	return NewCiPipelineYaml(c)
}

func NewDefaultCiPipeline() *CiPipeline {
	return &CiPipeline{
		Stages: wrapper.NewArray[string]().Push(stages.Enum().GetValues()...),
		Default: &CiPipelineDefault{
			AfterScript:  job.NewAfterScript(),
			BeforeScript: job.NewBeforeScript(),
		},
		Workflow: &CiPipelineWorkflow{
			Name: "default",
			Rules: &job.Rules{
				&job.Rule{
					If:   "$CI_MERGE_REQUEST_ID",
					When: job.NewWhen(enum.WhenAlways),
				},
				&job.Rule{
					If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
					When: job.NewWhen(enum.WhenAlways),
				},
				&job.Rule{
					If:   "$CI_MERGE_REQUEST_IID && $CI_PIPELINE_SOURCE == 'merge_request_event'",
					When: job.NewWhen(enum.WhenAlways),
				},
				&job.Rule{
					If:   "$CI_COMMIT_TAG != null",
					When: job.NewWhen(enum.WhenAlways),
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

func NewCiPipelineYaml(pipeline *CiPipeline) (*CiPipelineYaml, error) {
	err := pipeline.Validate()
	if err != nil {
		return nil, err
	}
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
	}, nil
}
func (c *CiPipelineYaml) AddToMap(parentMap map[string]interface{}) {
	parentMap["default"] = c.Default
	parentMap["workflow"] = c.Workflow
	parentMap["stages"] = c.Stages
	parentMap["variables"] = c.Variables
}

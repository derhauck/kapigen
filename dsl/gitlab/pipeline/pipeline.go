package pipeline

import (
	"gitlab.com/kateops/kapigen/dsl/enum"
	"gitlab.com/kateops/kapigen/dsl/gitlab/job"
	"gitlab.com/kateops/kapigen/dsl/gitlab/stages"
	"gitlab.com/kateops/kapigen/dsl/logger"
	"gitlab.com/kateops/kapigen/dsl/wrapper"
)

type CiPipelineDefault struct {
	AfterScript  wrapper.Array[string] `yaml:"after_script"`
	BeforeScript wrapper.Array[string] `yaml:"before_script"`
}

type CiPipelineWorkflow struct {
	Name  string     `yaml:"name"`
	Rules *job.Rules `yaml:"rules,omitempty"`
}

func (c *CiPipelineWorkflow) Validate() error {
	if c.Rules == nil {
		c.Rules = &job.Rules{}
	}

	if c.Name == "" {
		return wrapper.ErrorHandler("ci pipeline workflow name can not be empty", 2)
	}

	return nil
}

// CiPipeline is the root gitlab pipeline configuration
type CiPipeline struct {
	Stages       *wrapper.Array[string] `yaml:"stages"`
	Workflow     *CiPipelineWorkflow    `yaml:"workflow,omitempty"`
	AllowFailure job.AllowFailure       `yaml:"allow_failure,omitempty"`
	Default      *CiPipelineDefault     `yaml:"default,omitempty"`
	Variables    map[string]string      `yaml:"variables,omitempty"`
}

func (c *CiPipeline) Validate() error {
	err := c.Workflow.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (c *CiPipeline) DefaultCiPipeline() *CiPipeline {
	c.Stages = wrapper.NewArray[string]().Push(stages.Enum().GetValues()...)
	c.Workflow = &CiPipelineWorkflow{
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
	}

	c.Default = nil
	c.Variables = map[string]string{
		"KTC_STOP_PIPELINE": "false",
	}
	logger.DebugAny(c.Stages)
	return c
}

func (c *CiPipeline) Render() (*CiPipelineYaml, error) {
	return NewCiPipelineYaml(c)
}

type CiPipelineWorkflowYaml struct {
	Name  string                `yaml:"name"`
	Rules *job.RuleWorkflowYaml `yaml:"rules"`
}

func NewCiPipelineWorkflowYaml(workflow *CiPipelineWorkflow) *CiPipelineWorkflowYaml {
	if workflow == nil {
		return nil
	}
	return &CiPipelineWorkflowYaml{
		Name:  workflow.Name,
		Rules: workflow.Rules.GetRenderedWorkflowValue(),
	}
}

type CiPipelineDefaultYaml struct {
	AfterScript  []string `yaml:"after_script,omitempty"`
	BeforeScript []string `yaml:"before_script,omitempty"`
}

func NewCiPipelineDefaultYaml(defaultCiPipeline *CiPipelineDefault) *CiPipelineDefaultYaml {
	if defaultCiPipeline == nil {
		return nil
	}
	return &CiPipelineDefaultYaml{
		AfterScript:  defaultCiPipeline.AfterScript.Get(),
		BeforeScript: defaultCiPipeline.BeforeScript.Get(),
	}
}

type CiPipelineYaml struct {
	Default   *CiPipelineDefaultYaml  `yaml:"default,omitempty"`
	Workflow  *CiPipelineWorkflowYaml `yaml:"workflow,omitempty"`
	Stages    []string                `yaml:"stages"`
	Variables map[string]string       `yaml:"variables,omitempty"`
}

func NewCiPipelineYaml(pipeline *CiPipeline) (*CiPipelineYaml, error) {
	err := pipeline.Validate()
	if err != nil {
		return nil, err
	}
	return &CiPipelineYaml{
		Default:   NewCiPipelineDefaultYaml(pipeline.Default),
		Workflow:  NewCiPipelineWorkflowYaml(pipeline.Workflow),
		Stages:    wrapper.GetSlice(pipeline.Stages),
		Variables: pipeline.Variables,
	}, nil
}
func (c *CiPipelineYaml) AddToMap(parentMap map[string]interface{}) {
	parentMap["default"] = c.Default
	parentMap["workflow"] = c.Workflow
	parentMap["stages"] = c.Stages
	parentMap["variables"] = c.Variables
}

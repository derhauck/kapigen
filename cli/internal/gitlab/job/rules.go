package job

import (
	"errors"
	"fmt"

	"kapigen.kateops.com/internal/logger"
	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
)

type Rule struct {
	If           string              `yaml:"if,omitempty"`
	Changes      wrapper.StringSlice `yaml:"changes"`
	AllowFailure wrapper.Bool        `yaml:"allow_failure"`
	Variables    struct{}            `yaml:"variables"`
	When         WhenWrapper         `yaml:"when"`
}

func (r *Rule) AddChange(path string) *Rule {
	if !r.Changes.Has(path) {
		r.Changes.Add(path)
	}

	return r
}

type WhenWrapper struct {
	Value *when.When
}

func (w *WhenWrapper) Get() string {
	if w.Value == nil {
		return when.OnSuccess.String()
	}
	return w.Value.String()
}

func NewWhen(when when.When) WhenWrapper {
	return WhenWrapper{
		Value: &when,
	}
}

type Rules []*Rule

func (r *Rules) AddPathToChanges(path string) *Rules {
	for _, rule := range r.Get() {
		if len(rule.Changes.Get()) > 0 {
			rule.AddChange(path)
		}
	}
	return r
}

func (r *Rules) Get() []*Rule {
	return *r
}

func (r *Rules) Add(rule *Rule) *Rules {
	tmp := append(*r, rule)
	*r = tmp
	return r
}

func (r *Rules) AddRules(rules Rules) *Rules {
	for _, rule := range rules.Get() {
		r.Add(rule)
	}
	return r
}

type RuleYaml struct {
	If           string   `yaml:"if,omitempty"`
	Changes      []string `yaml:"changes,omitempty"`
	AllowFailure any      `yaml:"allow_failure,omitempty"`
	Variables    struct{} `yaml:"variables,omitempty"`
	When         string   `yaml:"when"`
}

type RuleWorkflowYaml []*RuleYaml

type RulesYaml []*RuleYaml

func (r *Rules) GetRenderedWorkflowValue() *RuleWorkflowYaml {
	return NewRulesWorkflowYaml(*r)
}
func (r *Rules) GetRenderedValue() *RulesYaml {
	return NewRulesYaml(*r)
}

func validateWorkflowRule(rule *Rule) error {
	if rule.If == "" {
		return errors.New("if is required")
	}

	if rule.When.Value != nil {
		whenValue := *rule.When.Value
		if whenValue.String() != when.Always.String() && whenValue.String() != when.Never.String() {
			return errors.New(fmt.Sprintf("when is not supported: %s", whenValue.String()))
		}
	}

	return nil
}
func NewRulesWorkflowYaml(rules Rules) *RuleWorkflowYaml {
	var rulesYaml = make(RuleWorkflowYaml, 0)
	for i := 0; i < len(rules); i++ {
		var currentRule = rules[i]
		if err := validateWorkflowRule(currentRule); err != nil {
			logger.Error(err.Error())
			continue
		}
		rulesYaml = append(rulesYaml, &RuleYaml{
			If:           currentRule.If,
			Changes:      currentRule.Changes.Get(),
			Variables:    currentRule.Variables,
			When:         currentRule.When.Get(),
			AllowFailure: nil,
		})
	}

	return &rulesYaml
}

func NewRulesYaml(rules Rules) *RulesYaml {
	var rulesYaml = make(RulesYaml, 0)
	for i := 0; i < len(rules); i++ {
		var currentRule = rules[i]
		rulesYaml = append(rulesYaml, &RuleYaml{
			If:           currentRule.If,
			Changes:      currentRule.Changes.Get(),
			AllowFailure: currentRule.AllowFailure.Get(),
			Variables:    currentRule.Variables,
			When:         currentRule.When.Get(),
		})
	}

	return &rulesYaml
}

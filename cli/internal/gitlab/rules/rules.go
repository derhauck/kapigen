package rules

import "kapigen.kateops.com/internal/pipeline/wrapper"

type Rule struct {
	If           string              `yaml:"if"`
	Changes      wrapper.StringSlice `yaml:"changes"`
	AllowFailure bool                `yaml:"allow_failure"`
	Variables    struct{}            `yaml:"variables"`
	When         WhenWrapper         `yaml:"when"`
}

type WhenWrapper struct {
	Value *WhenEnumType
}

func (w *WhenWrapper) Get() string {
	if w.Value == nil {
		return WhenEnumTypeOnSuccess.When()
	}
	return w.Value.When()
}

func NewWhen(when WhenEnumType) WhenWrapper {
	return WhenWrapper{
		Value: &when,
	}
}

type Rules []*Rule

func (r *Rules) AddPathToChanges(path string) *Rules {
	for _, rule := range r.Get() {
		if len(rule.Changes.Get()) > 0 {
			rule.Changes.Add(path)
		}
	}
	return r
}

func (r *Rules) Get() []*Rule {
	return *r
}

type WhenEnumType int

const (
	WhenEnumTypeOnSuccess WhenEnumType = iota
	WhenEnumTypeOnFailure
	WhenEnumTypeAlways
	WhenEnumTypeNever
)

func (c WhenEnumType) When() string {
	return []string{
		"on_success",
		"on_failure",
		"always",
		"never",
	}[c]
}

type RuleYaml struct {
	If           string   `yaml:"if"`
	Changes      []string `yaml:"changes,omitempty"`
	AllowFailure bool     `yaml:"allow_failure"`
	Variables    struct{} `yaml:"variables,omitempty"`
	When         string   `yaml:"when"`
}

type RulesYaml []*RuleYaml

func (r *Rules) GetRenderedValue() *RulesYaml {
	return NewRulesYaml(*r)
}
func NewRulesYaml(rules Rules) *RulesYaml {
	var rulesYaml = make(RulesYaml, 0)
	for i := 0; i < len(rules); i++ {
		var currentRule = rules[i]
		rulesYaml = append(rulesYaml, &RuleYaml{
			If:           currentRule.If,
			Changes:      currentRule.Changes.Get(),
			AllowFailure: currentRule.AllowFailure,
			Variables:    currentRule.Variables,
			When:         currentRule.When.Get(),
		})
	}

	return &rulesYaml
}

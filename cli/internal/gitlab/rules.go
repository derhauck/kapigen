package gitlab

import "kapigen.kateops.com/internal/pipeline/wrapper"

type Rule struct {
	If           string              `yaml:"if"`
	Changes      wrapper.StringSlice `yaml:"changes"`
	AllowFailure bool                `yaml:"allow_failure"`
	Variables    struct{}            `yaml:"variables"`
	When         WhenEnumType        `yaml:"when"`
}

type Rules []Rule

type WhenEnumType string

var WhenEnum = struct {
	OnSuccess WhenEnumType `yaml:"on_success"`
	OnFailure WhenEnumType `yaml:"on_failure"`
	Always    WhenEnumType `yaml:"always"`
}{
	OnSuccess: "success",
	OnFailure: "on_failure",
	Always:    "always",
}

type RuleYaml struct {
	If           string   `yaml:"if"`
	Changes      []string `yaml:"changes"`
	AllowFailure bool     `yaml:"allow_failure"`
	Variables    struct{} `yaml:"variables"`
	When         string   `yaml:"when"`
}

type RulesYaml []RuleYaml

type CiRules struct {
	Value Rules
}

func (c *CiRules) Get() Rules {
	return c.Value
}

func (c *CiRules) GetRenderedValue() RulesYaml {
	return NewRulesYaml(c.Value)
}
func NewRulesYaml(rules Rules) RulesYaml {
	var rulesYaml = make(RulesYaml, 0)
	for i := 0; i < len(rules); i++ {
		var currentRule = rules[i]
		rulesYaml = append(rulesYaml, RuleYaml{
			If:           currentRule.If,
			Changes:      currentRule.Changes.Get(),
			AllowFailure: currentRule.AllowFailure,
			Variables:    currentRule.Variables,
			When:         string(currentRule.When),
		})
	}

	return rulesYaml
}

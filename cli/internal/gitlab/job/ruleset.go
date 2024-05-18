package job

import (
	"fmt"

	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
)

func DefaultPipelineRules() *Rules {
	return &Rules{
		&Rule{
			If:   "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
			When: NewWhen(when.Never),
		},
		&Rule{
			If:      "$CI_MERGE_REQUEST_IID",
			Changes: *wrapper.NewStringSlice().Add("${KTC_PATH}/**/*"),
			When:    NewWhen(when.OnSuccess),
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
		&Rule{
			If:   "$CI_COMMIT_TAG != null",
			When: NewWhen(when.OnSuccess),
		},
	}
}

func DefaultOnlyReleasePipelineRules() *Rules {
	return &Rules{
		&Rule{
			If: "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
		},
		&Rule{
			If:   "$CI_COMMIT_TAG != null",
			When: NewWhen(when.OnSuccess),
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
	}
}

func RulesNotKapigen(rules *Rules) *Rules {
	for _, rule := range rules.Get() {
		rule.If = fmt.Sprintf("%s && %s", rule.If, "$CI_PROJECT_ID != \"57482547\"")
	}

	return rules
}

func RulesKapigen(rules *Rules) *Rules {
	for _, rule := range rules.Get() {
		rule.If = fmt.Sprintf("%s && %s", rule.If, "$CI_PROJECT_ID == \"57482547\"")
	}

	return rules
}

type DefaultPipelineRule struct {
	Changes wrapper.StringSlice
	Rules   *Rules
}

func (d *DefaultPipelineRule) GetChanges() wrapper.StringSlice {
	return d.Changes
}

func (d *DefaultPipelineRule) Get() *Rules {
	rules := *d.Rules
	for _, rule := range rules.Get() {
		if len(rule.Changes.Get()) > 0 {
			for _, change := range d.Changes.Get() {
				rule.AddChange(change)
			}
		}
	}
	return &rules
}

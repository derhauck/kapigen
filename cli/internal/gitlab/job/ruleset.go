package job

import (
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
			If:      "($CI_MERGE_REQUEST_IID || $CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH)",
			Changes: *wrapper.NewStringSlice().Add("${KTC_PATH}/**/*"),
			When:    NewWhen(when.Always),
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
	}
}

func DefaultReleasePipelineRules() *Rules {
	return &Rules{
		&Rule{
			If: "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
		},
		&Rule{
			If:      "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			Changes: *wrapper.NewStringSlice().Add("${KTC_PATH}/**/*"),
			When:    NewWhen(when.Always),
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
	}
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

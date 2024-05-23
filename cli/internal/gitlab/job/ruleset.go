package job

import (
	"fmt"
	"regexp"

	"kapigen.kateops.com/internal/pipeline/wrapper"
	"kapigen.kateops.com/internal/when"
)

func DefaultPipelineRules(paths []string) *Rules {
	rules := &Rules{}
	return rules.AddRules(*DefaultMergeRequestRules(paths)).
		AddRules(*DefaultMainBranchRules())
}
func DefaultMergeRequestRules(paths []string) *Rules {
	changes := wrapper.NewStringSlice()
	change := "**/*"
	for _, path := range paths {
		if path != "." {
			change = fmt.Sprintf("%s/%s", path, change)
		}
		r, _ := regexp.Compile("/+")
		if r != nil {
			change = r.ReplaceAllString(change, "/")
		}
		changes.Add(change)

	}
	return &Rules{
		&Rule{
			If:   "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
			When: NewWhen(when.Never),
		},
		&Rule{
			If:      "$CI_MERGE_REQUEST_IID",
			Changes: *changes,
			When:    NewWhen(when.OnSuccess),
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
	}
}

func DefaultMainBranchRules() *Rules {
	return &Rules{
		&Rule{
			If:   "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
			When: NewWhen(when.Never),
		},
		&Rule{
			If:   "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
			When: NewWhen(when.OnSuccess),
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
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

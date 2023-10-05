package rules

import "kapigen.kateops.com/internal/pipeline/wrapper"

func DefaultPipelineRules() *Rules {
	return &Rules{
		&Rule{
			If:   "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
			When: NewWhen(WhenEnumTypeNever),
		},
		&Rule{
			If: "($CI_MERGE_REQUEST_IID || $CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH) && $KTC_VERSION != \"0.0.0\"",
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
	}
}

func DefaultReleasePipelineRules() *Rules {
	return &Rules{
		&Rule{
			If:      "$KTC_STOP_PIPELINE != \"false\" && $DEBUG == null",
			Changes: *wrapper.NewStringSlice().Add("${KTC_PATH}/**/*"),
		},
		&Rule{
			If: "$CI_DEFAULT_BRANCH == $CI_COMMIT_BRANCH",
		},
		&Rule{
			If: "$KTC_TEST_PIPELINE",
		},
	}
}

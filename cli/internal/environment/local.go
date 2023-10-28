package environment

func SetLocalEnv() {
	CI_PROJECT_ID.SetIfEmpty("175")
	CI_MERGE_REQUEST_ID.SetIfEmpty("25")
	CI_MERGE_REQUEST_LABELS.SetIfEmpty("version::patch")
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.SetIfEmpty("feature/test")
	KAPIGEN_VERSION.SetIfEmpty("0.0.1")

}

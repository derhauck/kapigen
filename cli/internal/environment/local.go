package environment

func SetLocalEnv() {
	CI_PROJECT_ID.Set("175")
	CI_MERGE_REQUEST_ID.Set("25")
	CI_MERGE_REQUEST_LABELS.Set("version::patch")
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Set("feature/test")
	KAPIGEN_VERSION.Set("0.0.1")

}

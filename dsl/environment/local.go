package environment

func SetLocalEnv() {
	CI_PROJECT_ID.SetIfEmpty("57482547")
	CI_MERGE_REQUEST_ID.SetIfEmpty("25")
	CI_COMMIT_BRANCH.SetIfEmpty("feature/test")
	CI_DEFAULT_BRANCH.SetIfEmpty("main")
	CI_MERGE_REQUEST_LABELS.SetIfEmpty("version::patch")
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.SetIfEmpty("feature/test")
	KAPIGEN_VERSION.SetIfEmpty("0.0.1")
	CI_JOB_TOKEN.SetIfEmpty("invalid")
	CI_SERVER_URL.SetIfEmpty("https://gitlab.com")
	CI_PROJECT_DIR.SetIfEmpty("/app")
}

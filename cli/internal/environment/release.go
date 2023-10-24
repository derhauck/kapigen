package environment

func IsRelease() bool {
	_, err := CI_MERGE_REQUEST_ID.Lookup()
	if err != nil {
		return false
	}
	if CI_COMMIT_BRANCH.Get() != CI_DEFAULT_BRANCH.Get() {
		return false
	}
	return true
}

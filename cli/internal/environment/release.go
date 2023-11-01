package environment

import (
	"fmt"
	"kapigen.kateops.com/internal/logger"
	"regexp"
	"strconv"
)

func IsRelease() bool {
	_, err := CI_MERGE_REQUEST_ID.Lookup()
	if err == nil {
		return false
	}
	if CI_COMMIT_BRANCH.Get() != CI_DEFAULT_BRANCH.Get() {
		return false
	}
	return true
}

func GetBranchName() string {
	if IsRelease() {
		return CI_COMMIT_BRANCH.Get()
	}
	return CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Get()
}

func GetMergeRequestId() int {
	if IsRelease() {
		return getMergeRequestIdFromEnv()
	}
	return getMergeRequestIdFromCommit()
}

func getMergeRequestIdFromEnv() int {
	id, err := CI_MERGE_REQUEST_ID.Lookup()
	if err != nil {

	}
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		logger.Error(fmt.Sprintf("could not parse merge request id %s", id))
		return 0
	}
	return int(i)

}
func getMergeRequestIdFromCommit() int {
	message := CI_COMMIT_MESSAGE.Get()
	reg := regexp.MustCompile("![0-9]*$")
	stringId := reg.FindString(message)
	if stringId == "" {
		logger.Error("no merge request id found in commit message")
		return 0
	}
	intId, err := strconv.ParseInt(stringId, 10, 32)
	if err != nil {
		logger.ErrorE(err)
		return 0
	}
	return int(intId)
}

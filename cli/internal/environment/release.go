package environment

import (
	"fmt"
	"regexp"
	"strconv"

	"kapigen.kateops.com/internal/logger"
)

func IsRelease() bool {
	commit, _ := CI_COMMIT_BRANCH.Lookup()
	def, _ := CI_DEFAULT_BRANCH.Lookup()
	_, errTag := CI_COMMIT_TAG.Lookup()
	if commit == def && errTag == nil {
		return true
	} else {
		return false
	}
}

func GetBranchName() string {
	if IsRelease() {
		return CI_COMMIT_BRANCH.Get()
	}
	return CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Get()
}

func GetMergeRequestId() int {
	if IsRelease() {
		return getMergeRequestIdFromCommit(CI_COMMIT_MESSAGE.Get())
	}
	return getMergeRequestIdFromEnv()
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
func getMergeRequestIdFromCommit(message string) int {
	reg := regexp.MustCompile("![0-9]+")
	stringId := reg.FindString(message)
	if stringId == "" {
		logger.Error(fmt.Sprintf("no merge request id found in commit message: %s", message))
		return 0
	}
	intId, err := strconv.ParseInt(stringId[1:], 10, 32)
	if err != nil {
		logger.ErrorE(err)
		return 0
	}
	return int(intId)
}

package environment

import (
	"regexp"
	"strconv"

	"gitlab.com/kateops/kapigen/dsl/logger"
)

func IsRelease() bool {
	_, errTag := CI_COMMIT_TAG.Lookup()
	if errTag == nil {
		return true
	} else {
		return false
	}
}

func GetBranchName() string {
	if _, err := CI_MERGE_REQUEST_ID.Lookup(); err != nil {
		return CI_COMMIT_BRANCH.Get()
	}
	return CI_MERGE_REQUEST_SOURCE_BRANCH_NAME.Get()
}

func GetMergeRequestId() int {
	if _, err := CI_MERGE_REQUEST_ID.Lookup(); err != nil {
		return getMergeRequestIdFromCommit(CI_COMMIT_MESSAGE.Get())
	}
	return getMergeRequestIdFromEnv()
}

func getMergeRequestIdFromEnv() int {
	id := CI_MERGE_REQUEST_ID.Get()
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		logger.Errorf("could not parse merge request id %s", id)
		return 0
	}
	return int(i)

}
func getMergeRequestIdFromCommit(message string) int {
	reg := regexp.MustCompile("![0-9]+")
	stringId := reg.FindString(message)
	if stringId == "" {
		logger.Errorf("no merge request id found in commit message: %s", message)
		return 0
	}
	intId, err := strconv.ParseInt(stringId[1:], 10, 32)
	if err != nil {
		logger.ErrorE(err)
		return 0
	}
	return int(intId)
}

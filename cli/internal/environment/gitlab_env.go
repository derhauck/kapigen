package environment

import (
	"errors"
	"fmt"
	"os"

	"kapigen.kateops.com/internal/logger"
)

type Variable int

const (
	CI_PROJECT_ID Variable = iota
	CI_MERGE_REQUEST_ID
	CI_COMMIT_BRANCH
	CI_COMMIT_MESSAGE
	CI_COMMIT_SHA
	CI_COMMIT_TAG
	CI_DEFAULT_BRANCH
	CI_JOB_TOKEN
	CI_MERGE_REQUEST_LABELS
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME
	CI_PIPELINE_TOKEN
	CI_SERVER_HOST
	GITLAB_TOKEN
	LOS_AUTH_TOKEN
	KAPIGEN_VERSION
)

var values = map[Variable]string{
	CI_COMMIT_SHA:                       "CI_COMMIT_SHA",
	CI_COMMIT_BRANCH:                    "CI_COMMIT_BRANCH",
	CI_COMMIT_MESSAGE:                   "CI_COMMIT_MESSAGE",
	CI_COMMIT_TAG:                       "CI_COMMIT_TAG",
	CI_DEFAULT_BRANCH:                   "CI_DEFAULT_BRANCH",
	CI_JOB_TOKEN:                        "CI_JOB_TOKEN",
	CI_PROJECT_ID:                       "CI_PROJECT_ID",
	CI_PIPELINE_TOKEN:                   "CI_PIPELINE_TOKEN",
	CI_MERGE_REQUEST_ID:                 "CI_MERGE_REQUEST_ID",
	CI_MERGE_REQUEST_LABELS:             "CI_MERGE_REQUEST_LABELS",
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME: "CI_MERGE_REQUEST_SOURCE_BRANCH_NAME",
	KAPIGEN_VERSION:                     "KAPIGEN_VERSION",
	GITLAB_TOKEN:                        "GITLAB_TOKEN",
	LOS_AUTH_TOKEN:                      "LOS_AUTH_TOKEN",
	CI_SERVER_HOST:                      "CI_SERVER_HOST",
}

func (v Variable) Key() string {
	if v, ok := values[v]; ok {
		return v
	}
	logger.Error(fmt.Sprintf("not found env var for id: '%d'", v))
	return ""
}

func (v Variable) Set(value string) {
	if key, ok := values[v]; ok {
		err := os.Setenv(key, value)
		if err != nil {
			logger.ErrorE(err)
		}
	}
}

func (v Variable) SetIfEmpty(value string) bool {
	if _, err := v.Lookup(); err != nil {
		v.Set(value)
		return true
	}
	return false
}

func (v Variable) Get() string {
	value := os.Getenv(v.Key())
	if value == "" {
		logger.Error(fmt.Sprintf("env var '%s' is not set", v.Key()))
	}
	return value
}

func (v Variable) Lookup() (string, error) {
	value, set := os.LookupEnv(v.Key())
	if set {
		return value, nil
	}
	return value, errors.New(fmt.Sprintf("env var '%s' is not set", v.Key()))
}

func (v Variable) Unset() {
	err := os.Unsetenv(v.Key())
	if err != nil {
		logger.ErrorE(err)
	}
}

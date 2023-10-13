package environment

import (
	"errors"
	"fmt"
	"kapigen.kateops.com/internal/logger"
	"os"
)

type Variable int

const (
	CI_PROJECT_ID Variable = iota
	CI_MERGE_REQUEST_ID
	CI_COMMIT_SHA
	KAPIGEN_VERSION
	CI_MERGE_REQUEST_LABELS
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME
	LOS_AUTH_TOKEN
)

var values = map[Variable]string{
	CI_COMMIT_SHA:                       "CI_COMMIT_SHA",
	CI_PROJECT_ID:                       "CI_PROJECT_ID",
	CI_MERGE_REQUEST_ID:                 "CI_MERGE_REQUEST_ID",
	CI_MERGE_REQUEST_LABELS:             "CI_MERGE_REQUEST_LABELS",
	CI_MERGE_REQUEST_SOURCE_BRANCH_NAME: "CI_MERGE_REQUEST_SOURCE_BRANCH_NAME",
	KAPIGEN_VERSION:                     "KAPIGEN_VERSION",
	LOS_AUTH_TOKEN:                      "LOS_AUTH_TOKEN",
}

func (v Variable) String() string {
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

func (v Variable) Get() string {
	value := os.Getenv(v.String())
	if value == "" {
		logger.Error(fmt.Sprintf("env var '%s' is not set", v.String()))
	}
	return value
}
func (v Variable) Lookup() (string, error) {
	value, set := os.LookupEnv(v.String())
	if set {
		return value, nil
	}
	return value, errors.New(fmt.Sprintf("env var '%s' is not set", v.String()))
}

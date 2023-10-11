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
)

var values = map[Variable]string{
	CI_PROJECT_ID:       "CI_PROJECT_ID",
	CI_MERGE_REQUEST_ID: "CI_MERGE_REQUEST_ID",
}

func (v Variable) String() string {
	if v, ok := values[v]; ok {
		return v
	}
	logger.Error(fmt.Sprintf("not found env var for id: '%d'", v))
	return ""
}

func Get(key Variable) string {
	name := values[key]
	value := os.Getenv(name)
	return value
}
func Lookup(key Variable) (string, error) {
	name := values[key]
	value, set := os.LookupEnv(name)
	if set {
		return value, nil
	}
	return value, errors.New(fmt.Sprintf("env var '%s' is not set", name))
}

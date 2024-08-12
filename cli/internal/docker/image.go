package docker

import (
	"fmt"

	"gitlab.com/kateops/kapigen/dsl/logger"
)

type Image int

var DEPENDENCY_PROXY = ""

const (
	Kapigen_Latest Image = iota
	Alpine_3_18
	Terraform_Base
	BUILDKIT
	BUILDKIT_ROTLESS
	CRANE_DEBUG

	GOLANG_1_21
	GOLANG_GOLANGCI_LINT
)

var values = map[Image]struct {
	Image string
	Proxy bool
}{
	Kapigen_Latest:       {"registry.gitlab.com/kateops/kapigen/cli:latest", false},
	Alpine_3_18:          {"alpine:3.18", true},
	Terraform_Base:       {"hub.kateops.com/base/terraform:latest", false},
	BUILDKIT:             {"moby/buildkit:master", true},
	BUILDKIT_ROTLESS:     {"moby/buildkit:v0.12.3-rootless", true},
	CRANE_DEBUG:          {"gcr.io/go-containerregistry/crane:debug", false},
	GOLANG_1_21:          {"golang:1.21", true},
	GOLANG_GOLANGCI_LINT: {"golangci/golangci-lint:v1.59.1", true},
}

func (c Image) String() string {
	if value, ok := values[c]; ok {
		if value.Proxy {
			return fmt.Sprintf("%s%s", DEPENDENCY_PROXY, value.Image)
		}

		return value.Image
	}
	logger.Errorf("Not found for id: '%d'", c)

	if values[Alpine_3_18].Proxy {
		return fmt.Sprintf("%s%s", DEPENDENCY_PROXY, values[Alpine_3_18].Image)
	}
	return values[Alpine_3_18].Image
}

GO_IMAGE=golang:1.21
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))
DOCKER_RUN=docker run --rm $(DOCKER_ARGS) -v $(mkfile_dir):/app -u $(shell id -u):$(shell id -g)

.PHONY: cli
cli-go: DOCKER_ARGS=-it -e GOCACHE=/app/.cache
cli-go:
	$(DOCKER_RUN) $(GO_IMAGE) bash
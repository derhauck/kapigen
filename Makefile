NPM_IMAGE=node:current-alpine3.15
GO_IMAGE=golang:1.21rc3
DOCKER_RUN=docker run --rm $(DOCKER_ARGS) -v ${PWD}:/app -w /app -u $(shell id -u):$(shell id -g)
.PHONY: cli
cli: DOCKER_ARGS=-it
cli:
	$(DOCKER_RUN) $(NPM_IMAGE) sh

.PHONY: cli-go
cli-go: DOCKER_ARGS=-it -e GOCACHE=/app/.cache
cli-go:
	$(DOCKER_RUN) $(GO_IMAGE) bash
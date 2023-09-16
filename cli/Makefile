include ../Makefile
GO_SRC_PATH=/app
GO_ARGS=-it -e GOMODCACHE=/app/.pkg -e GOCACHE=/app/.cache -e ENVIRONMENT=DEV
.PHONY: cli
cli: DOCKER_ARGS=$(GO_ARGS)
cli:
	$(DOCKER_RUN) $(GO_IMAGE) bash

.PHONY: run
run: DOCKER_ARGS=$(GO_ARGS)
run:
	$(DOCKER_RUN) -w $(GO_SRC_PATH) $(GO_IMAGE) go run .

.PHONH: test
test: DOCKER_ARGS=$(GO_ARGS)
test:
	$(DOCKER_RUN)  -w $(GO_SRC_PATH) $(GO_IMAGE) go test ./...
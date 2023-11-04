include ../Makefile
GO_SRC_PATH=/app
GO_ARGS=-it -e GOMODCACHE=/app/.pkg -e GOCACHE=/app/.cache -e ENVIRONMENT=DEV --env-file .env
.PHONY: cli
cli: DOCKER_ARGS=$(GO_ARGS)
cli:
	$(DOCKER_RUN) $(GO_IMAGE) bash

.PHONY: run
run: DOCKER_ARGS=$(GO_ARGS)
run:
	$(DOCKER_RUN) -w $(GO_SRC_PATH) $(GO_IMAGE) go run . $(ARGS)

.PHONY: build
build: DOCKER_ARGS=$(GO_ARGS)
build:
	$(DOCKER_RUN) -w $(GO_SRC_PATH) $(GO_IMAGE) go build -o .bin/kapigen .
	chmod +x bin/kapigen

.PHONY: test
test: DOCKER_ARGS=$(GO_ARGS)
test:
	$(DOCKER_RUN)  -w $(GO_SRC_PATH) $(GO_IMAGE) go test ./...  -coverpkg=kapigen.kateops.com/internal/... -coverprofile=profile.cov
	$(DOCKER_RUN)  -w $(GO_SRC_PATH) $(GO_IMAGE) go tool cover -func profile.cov

.PHONY: pipeline
pipeline:
	$(MAKE) run ARGS="generate pipeline -v"
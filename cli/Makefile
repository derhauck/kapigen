include ../Makefile
GO_SRC_PATH=/app/cli
DOCKER_ARGS=-it -e GOMODCACHE=/app/cli/.pkg -e GOCACHE=/app/cli/.cache -e ENVIRONMENT=DEV -w /app/cli

.PHONY: cli
cli: DOCKER_ARGS=$(GO_ARGS)
cli:
	$(DOCKER_RUN) $(GO_IMAGE) bash

.PHONY: run
run:
	$(DOCKER_RUN) -w $(GO_SRC_PATH) $(GO_IMAGE) go run . $(ARGS)

.PHONY: build
build:
	$(DOCKER_RUN) -w $(GO_SRC_PATH) $(GO_IMAGE) go build -o .bin/kapigen .
	chmod +x bin/kapigen

.PHONY: test
test:
	$(DOCKER_RUN)  -w $(GO_SRC_PATH) $(GO_IMAGE) go test ./... -v  -coverpkg=gitlab.com/kateops/kapigen/cli/internal/... -coverprofile=profile.cov
	$(DOCKER_RUN)  -w $(GO_SRC_PATH) $(GO_IMAGE) go tool cover -func profile.cov

.PHONY: test-update-snapshots
test-update-snapshots:
	$(DOCKER_RUN)  -e  UPDATE_SNAPS=true -w $(GO_SRC_PATH) $(GO_IMAGE) go test ./... -coverpkg=gitlab.com/kateops/kapigen/cli/internal/... -v


.PHONY: test-report
test-report:
	$(DOCKER_RUN) --entrypoint=/bin/sh -w $(GO_SRC_PATH) $(GO_IMAGE) -c "go-junit-report -parser gojson -iocopy -out junit.xml"

.PHONY: pipeline
pipeline:
	$(MAKE) run ARGS="pipeline generate -v --mode gitlab"

.PHONY: test-pipeline
test-pipeline:
	$(MAKE) run ARGS="pipeline generate -v --config test.kapigen.yaml"

.PHONY: test-pipeline-no-merge
test-pipeline-no-merge:
	$(MAKE) run ARGS="pipeline generate -v --config test.kapigen.yaml --no-merge --file pipeline-no-merge.yaml"

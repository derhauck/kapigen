include ../Makefile
GO_SRC_PATH=/app/cli
GO_ARGS=-it -e GOMODCACHE=/app/cli/.pkg -e GOCACHE=/app/cli/.cache -e ENVIRONMENT=DEV
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
.PHONY: test-report
test-report: DOCKER_ARGS=$(GO_ARGS)
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

.PHONY: lint
lint:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run


.PHONY: lint-fix
lint-fix:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run -v --fix


.PHONY: lint-report
lint-report:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golaps ngci-lint:v1.59.1 golangci-lint run -v --out-format=junit-xml:junit.xml
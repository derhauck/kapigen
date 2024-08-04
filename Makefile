GO_IMAGE=golang:1.21
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
mkfile_dir := $(dir $(mkfile_path))
DOCKER_RUN=docker run --rm $(DOCKER_ARGS) -v $(mkfile_dir):/app -u $(shell id -u):$(shell id -g)

.PHONY: cli
cli-go: DOCKER_ARGS=-it -e GOCACHE=/app/.cache
cli-go:
	$(DOCKER_RUN) $(GO_IMAGE) bash


.PHONY: lint
lint:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run


.PHONY: lint-fix
lint-fix:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run -v --fix


.PHONY: lint-report
lint-report:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golaps ngci-lint:v1.59.1 golangci-lint run -v --out-format=junit-xml:junit.xml

.PHONY: lcc # line count complete (including tests)
lcc: params=. -name "*.go" -not -path "**/.pkg/*"
lcc: -lc

.PHONY: lc # line count
lc: params=. -name "*.go" -not -path "**/.pkg/*" -not -name "*_test.go"
lc: -lc

.PHONY: -lc
-lc:
	@echo "Files: $(params)"
	@echo "================================================="
	@find $(params)
	@echo "================================================="
	@echo "Total:"
	@find $(params) -exec cat {} \; | wc -l
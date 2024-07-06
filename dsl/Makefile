include ../Makefile

GO_SRC_PATH=/app/dsl
DOCKER_ARGS=-it -e GOMODCACHE=/app/dsl/.pkg -e GOCACHE=/app/dsl/.cache -w $(GO_SRC_PATH)

.PHONY: cli
cli:
	$(DOCKER_RUN) $(GO_IMAGE) bash

.PHONY: test
test:
	$(DOCKER_RUN) $(GO_IMAGE) go test ./...  -coverpkg=gitlab.com/kateops/kapigen/dsl/... -coverprofile=profile.cov
	$(DOCKER_RUN) $(GO_IMAGE) go tool cover -func profile.cov

.PHONY: lint
lint:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run


.PHONY: lint-fix
lint-fix:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run -v --fix


.PHONY: lint-report
lint-report:
	$(DOCKER_RUN) -v ${PWD}/.cache:/.cache golangci/golangci-lint:v1.59.1 golangci-lint run -v --out-format=junit-xml:junit.xml
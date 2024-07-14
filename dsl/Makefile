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

.PHONY: test-update-snapshots
test-update-snapshots:
	$(DOCKER_RUN)  -e  UPDATE_SNAPS=true -w $(GO_SRC_PATH) $(GO_IMAGE) go test ./... -coverpkg=gitlab.com/kateops/kapigen/dsl/... -v


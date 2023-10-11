FROM gitlab.kateops.com/infrastructure/dependency_proxy/containers/golang:1.21.3-alpine3.18 as builder
#ENV CGO_ENABLED=0

ENV GOMODCACHE=/app/.pkg
ENV GOCACHE=/app/.cache

COPY cmd /app/cmd
COPY internal /app/internal
COPY go.mod go.sum main.go /app/
#COPY .cache /app/.cache
#COPY .pkg /app/.pkg
WORKDIR /app

RUN --mount=type=cache,target=/app/.cache \
    --mount=type=cache,target=/app/.pkg \
    go mod download && \
    go build -o kapigen
#    go build -tags netgo -a -v -o kapigen
FROM gitlab.kateops.com/infrastructure/dependency_proxy/containers/alpine:3.18.4
RUN apk add --no-cache libc6-compat
COPY --from=builder /app/kapigen /usr/bin/kapigen
ENTRYPOINT ["/usr/bin/kapigen"]
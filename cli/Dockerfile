FROM gitlab.kateops.com/infrastructure/dependency_proxy/containers/golang:1.21rc3 as builder
COPY . /app
RUN cd /app && \
    go build -o kapigen .
FROM gitlab.kateops.com/infrastructure/dependency_proxy/containers/alpine:3.18.4
COPY --from=builder /app/kapigen /usr/bin/kapigen

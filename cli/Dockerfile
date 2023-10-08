FROM golang:1.21rc3 as builder
COPY . /app
RUN cd /app && \
    go build -o kapigen .
FROM alpine:3.18.4
COPY --from=builder /app/kapigen /usr/bin/kapigen

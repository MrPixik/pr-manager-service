FROM golang:1.24 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg ./pkg
COPY internal ./internal/
COPY cmd ./cmd/

RUN go build -v -o pr-manager-service ./cmd/main.go

FROM debian:stable-slim

WORKDIR /app

COPY --from=builder /usr/src/app/pr-manager-service .
CMD ["./pr-manager-service"]
FROM golang:1.21 AS build

WORKDIR /build/gentree

COPY go.mod go.sum ./
RUN go mod download

RUN ["go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest"]

COPY gentree/*.go ./
ENTRYPOINT ["golangci-lint", "run"]

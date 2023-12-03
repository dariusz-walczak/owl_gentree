FROM golang:1.21 AS build

WORKDIR /build/gentree

COPY go.mod go.sum ./
RUN go mod download

RUN ["go", "install", "github.com/rakyll/gotest@latest"]

ENV GIN_MODE=release

COPY gentree/*.go ./
COPY run-unit-tests.sh ./

ENTRYPOINT ["./run-unit-tests.sh"]

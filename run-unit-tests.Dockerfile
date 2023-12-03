FROM golang:1.21 AS build

WORKDIR /build/gentree

COPY go.mod go.sum ./
RUN go mod download

RUN ["go", "install", "github.com/rakyll/gotest@latest"]

ENV GIN_MODE=release

COPY gentree/*.go ./
COPY run-unit-tests.sh ./

# Create an empty output directory as a fallback for the docker run without the output directory
#  mounted (as it happens in case of the github action run):
RUN ["mkdir", "/output"]

ENTRYPOINT ["./run-unit-tests.sh"]

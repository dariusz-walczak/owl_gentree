FROM golang AS build

WORKDIR /build/gentree
COPY gentree/go.mod .
RUN ["go", "get", "-u", "github.com/gin-gonic/gin"]
RUN ["go", "get", "-u", "github.com/sirupsen/logrus"]
RUN ["go", "get", "-u", "github.com/jessevdk/go-flags"]
RUN ["go", "get", "-u", "github.com/gin-contrib/location"]
RUN ["go", "get", "-u", "github.com/stretchr/testify/assert"]
RUN ["go", "get", "-u", "github.com/stretchr/testify/require"]
RUN ["go", "install", "github.com/rakyll/gotest@latest"]

ENV GIN_MODE=release

COPY gentree/*.go ./
ENTRYPOINT ["gotest", ".", "-v", "-cover"]

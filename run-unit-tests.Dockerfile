FROM golang AS build

WORKDIR /build/gentree
COPY gentree/go.mod .
RUN ["go", "get", "-u", "github.com/gin-gonic/gin"]
RUN ["go", "get", "-u", "github.com/sirupsen/logrus"]
RUN ["go", "get", "-u", "github.com/jessevdk/go-flags"]
RUN ["go", "get", "-u", "github.com/gin-contrib/location"]
RUN ["go", "get", "-u", "github.com/stretchr/testify/assert"]
RUN ["go", "get", "-u", "github.com/stretchr/testify/require"]

ENV GIN_MODE=release

COPY gentree/*.go ./
COPY run_gentree_ut.sh ./
ENTRYPOINT ["./run_gentree_ut.sh"]
#ENTRYPOINT ["go", "test", "-v", "-cover", "."]

FROM golang AS build

WORKDIR /build/gentree
COPY gentree/go.mod .
RUN ["go", "get", "-u", "github.com/gin-gonic/gin"]
RUN ["go", "get", "-u", "github.com/sirupsen/logrus"]
RUN ["go", "get", "-u", "github.com/jessevdk/go-flags"]
RUN ["go", "get", "-u", "github.com/gin-contrib/location"]

COPY  gentree/*.go ./
RUN CGO_ENABLED=0 go build


FROM alpine

EXPOSE 8080/tcp

COPY --from=build /build/gentree/gentree /app/gentree
WORKDIR /app
ENTRYPOINT ["./gentree"]

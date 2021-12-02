FROM golang AS build

WORKDIR /build/gentree
COPY gentree/go.mod .
RUN ["go", "get", "-u", "github.com/gin-gonic/gin"]
RUN ["go", "get", "-u", "github.com/sirupsen/logrus"]

COPY  gentree/*.go .
RUN CGO_ENABLED=0 go build


FROM golang

EXPOSE 8080/tcp

COPY --from=build /build/gentree/gentree /app/gentree
WORKDIR /app
CMD ["./gentree"]

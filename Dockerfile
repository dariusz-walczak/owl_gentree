FROM golang:1.21 AS build

WORKDIR /build/gentree

COPY go.mod go.sum ./
RUN go mod download

COPY  gentree/*.go ./
RUN CGO_ENABLED=0 go build

FROM alpine

EXPOSE 8080/tcp

COPY --from=build /build/gentree/gentree /app/gentree
WORKDIR /app
ENTRYPOINT ["./gentree"]

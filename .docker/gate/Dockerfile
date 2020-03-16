ARG GO_VER

FROM golang:${PHP_VER}-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -o ./builds/linux/gate ./cmd/gate.go

FROM alpine:3.7

COPY --from=builder /app/builds/linux/gate /usr/bin/gate
RUN addgroup minteruser && adduser -D -h /minter -G minteruser minteruser
USER minteruser
WORKDIR /minter
CMD ["/usr/bin/gate"]
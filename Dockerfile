FROM golang:1.13.7-alpine as builder

WORKDIR /app
COPY ./ /app
RUN apk add --no-cache make gcc musl-dev linux-headers
RUN go mod download
RUN go build -o ./builds/linux/gate ./cmd/gate.go

FROM golang:1.13.7-alpine

COPY --from=builder /app/builds/linux/gate /usr/bin/gate
ENTRYPOINT ["/usr/bin/gate"]
CMD ["gate"]

FROM golang:1.13.7-alpine

RUN apk add --no-cache make gcc musl-dev linux-headers
RUN go get -u github.com/go-delve/delve/cmd/dlv
WORKDIR /app

COPY .docker/go/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["dlv", "debug", "/app/cmd/gate.go", "--headless", "--listen=:2345", "--api-version=2", "--log"]

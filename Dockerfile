FROM golang:latest as builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build

######## Start a new stage from scratch #######
FROM alpine:latest

WORKDIR /root/

RUN apk update && apk upgrade && apk --no-cache add dnsmasq

COPY --from=builder /app/ext-dnsmasq .

CMD ["./ext-dnsmasq", "-h"]

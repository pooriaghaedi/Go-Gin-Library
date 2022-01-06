# syntax=docker/dockerfile:1

FROM golang:1.17-alpine
RUN export GOPATH="$(pwd)/config:$(pwd)"
RUN apk add git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /docker-go-gin

EXPOSE 8080

CMD [ "/docker-go-gin" ]
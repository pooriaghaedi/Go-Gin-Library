# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

# RUN apk add git
WORKDIR /app
ENV GOPATH="/app/config:/app"
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY ./ /app

RUN go build -o /docker-go-gin

EXPOSE 8080

CMD [ "/docker-go-gin" ]
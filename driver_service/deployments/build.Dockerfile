FROM golang:1.21.1-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY app ./app
COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY ./.env ./.env

WORKDIR /app/cmd
RUN go build -o app

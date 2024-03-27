FROM golang:1.21.1-alpine

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY internal ./internal
COPY cmd ./cmd
COPY ./.env.dev ./.env.dev

WORKDIR /app/cmd
RUN go build -o app


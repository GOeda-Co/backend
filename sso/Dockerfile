FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

COPY . .

ENV CONFIG_PATH=./config/config.yaml

RUN go build -o server cmd/sso/main.go

FROM alpine:latest

RUN adduser -D appuser

WORKDIR /app

COPY --from=builder /app/server .

COPY --from=builder /app/config ./config

ENV CONFIG_PATH=./config/config.yaml

USER appuser

EXPOSE 44044

ENTRYPOINT ["./server"]


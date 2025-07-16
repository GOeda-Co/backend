# build-stage
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .
RUN go build -o server cmd/stats/main.go

# production-stage
FROM alpine:latest
RUN adduser -D appuser
WORKDIR /app

COPY --from=builder /app/server .

COPY ./config ./config
COPY .env .env

USER appuser
EXPOSE 8301
ENTRYPOINT ["./server"]
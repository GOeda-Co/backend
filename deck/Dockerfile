# build-stage
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o server cmd/deck/main.go

# production-stage
FROM alpine:latest
RUN adduser -D appuser
WORKDIR /app
COPY --from=builder /app/server .
USER appuser
EXPOSE 8202
ENTRYPOINT ["/app/server"]

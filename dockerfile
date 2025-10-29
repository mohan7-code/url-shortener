# Stage 1: Build
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o url-shortener ./main.go

# Stage 2: Run
FROM alpine:latest
WORKDIR /app

# Copy app binary and goose binary
COPY --from=builder /app/url-shortener .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/migrations ./migrations
COPY .env .

EXPOSE 8080
CMD ["./url-shortener"]

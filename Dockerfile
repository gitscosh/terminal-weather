FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install certificates for HTTPS calls
RUN apk add --no-cache ca-certificates

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o terminal-weather

FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Copy compiled binary from builder stage
COPY --from=builder /app/terminal-weather /usr/bin/weather

ENTRYPOINT ["weather"]
CMD ["--help"]
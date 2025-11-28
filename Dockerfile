# Multi-stage Docker build for hulud-scan
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty)" \
    -o hulud-scan

# Final stage - minimal runtime image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/hulud-scan /usr/local/bin/hulud-scan

# Create non-root user
RUN addgroup -g 1000 hulud && \
    adduser -D -u 1000 -G hulud hulud

USER hulud

# Default command
ENTRYPOINT ["hulud-scan"]
CMD ["--help"]

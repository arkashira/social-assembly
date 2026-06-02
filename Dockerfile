# Dockerfile
# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies for building
RUN apk add --no-cache ca-certificates git gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /social-assembly

# Runtime stage
FROM alpine:3.19

# Install CA certificates and dumb-init for proper signal handling
RUN apk add --no-cache ca-certificates dumb-init

# Create non-root user
RUN adduser -D appuser
USER appuser
WORKDIR /home/appuser

# Copy binary from builder
COPY --from=builder --chown=appuser:appuser /social-assembly /usr/local/bin/social-assembly

# Copy static files if any
COPY --chown=appuser:appuser static/ ./static/ || true

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Entrypoint with dumb-init for proper signal handling
ENTRYPOINT ["dumb-init", "--"]

# Command to run the application
CMD ["social-assembly"]
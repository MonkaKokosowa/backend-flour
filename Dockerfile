# ---------- Build Stage ----------
FROM golang:1.24-alpine AS builder

# Install git if you need to pull private modules
RUN apk add --no-cache git

# Set working directory inside the container
WORKDIR /app

# Cache go mod and dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary with stripped debug info for smaller size
RUN go build -ldflags="-s -w" -o backend-flour ./cmd/server.go

# ---------- Final Stage ----------
FROM alpine:latest

# Install ca-certificates (needed for HTTPS, email, etc)
RUN apk add --no-cache ca-certificates

# Create a non-root user (optional but recommended for security)
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Copy the compiled binary from builder
COPY --from=builder /app/backend-flour .

# Set permissions: switch to non-root user
USER appuser

# Run the binary
CMD ["./backend-flour"]

LABEL org.opencontainers.image.source="https://github.com/MonkaKokosowa/backend-flour"
LABEL org.opencontainers.image.description="This image is responsible for handling backend of my personal website (wheatflour.pl). As of writing this it's supporting sending mail and in future it will support blog posting and commenting capabilities."
LABEL org.opencontainers.image.licenses="GPL-3.0-only"

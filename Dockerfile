# Builder stage
FROM golang:1.25-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=0 creates a statically linked binary
# -ldflags="-w -s" reduces binary size by removing debug information
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/parallax-game .

# Final stage
FROM alpine:latest

LABEL org.opencontainers.image.source="https://github.com/snowlynxsoftware/parallax-game"

# Install minimal runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    update-ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/parallax-game /app/

# Expose the application port
EXPOSE 3000

# Run the application
CMD ["/app/parallax-game"]
# Multi-stage build for optimized Docker image

# Stage 1: Build stage
FROM golang:1.25.0-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o corestack cmd/server/main.go

# Stage 2: Runtime stage
FROM alpine:latest

# Install runtime dependencies (for timezone support and CA certificates)
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/corestack .

# Create .env file location (can be mounted or provided via environment variables)
RUN mkdir -p /root/config

# Expose port
EXPOSE 5000

# Run the application
CMD ["./corestack"]

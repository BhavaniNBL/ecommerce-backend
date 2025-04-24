# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required tools
RUN apk add --no-cache git

# Copy go.mod and go.sum from project root to leverage dependency caching
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy the Payment Service source code
COPY . .

# Move to payment-service directory
WORKDIR /app/services/payment-service

# Build the Payment Service binary
RUN go build -o payment-service main.go

# ----------- Final stage ----------- 
FROM alpine:3.18

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/services/payment-service/payment-service .

# Expose HTTP port (adjust if using gRPC as well)
EXPOSE 8084

# Start the service
CMD ["./payment-service"]

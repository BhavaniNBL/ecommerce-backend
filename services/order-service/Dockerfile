# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required tools
RUN apk add --no-cache git

# Copy go.mod and go.sum to leverage caching
COPY ../../go.mod ../../go.sum ./
RUN go mod download

# Copy the Order Service source code
COPY . ./

# Move to order-service directory
WORKDIR /app/services/order-service

# Build the Order Service binary
RUN go build -o order-service main.go

# ----------- Final stage -----------
FROM alpine:3.18

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/services/order-service/order-service .

# Expose HTTP port
EXPOSE 8081

CMD ["./order-service"]

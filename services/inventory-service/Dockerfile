# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY . ./

# Move to inventory-service source directory
WORKDIR /app/services/inventory-service

# Build the binary
RUN go build -o inventory-service main.go

# ----------- Final stage -----------
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/services/inventory-service/inventory-service .

EXPOSE 50053
EXPOSE 8083

CMD ["./inventory-service"]

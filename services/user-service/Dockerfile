# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install required packages
RUN apk add --no-cache git

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY . .

# Move to user-service source directory
WORKDIR /app/services/user-service

# Build the binary
RUN go build -o user-service main.go

# RUN go build -o user-service main.go
# ----------- Final stage -----------
FROM alpine:3.18

WORKDIR /app
    
# Copy binary from builder
COPY --from=builder /app/services/user-service/user-service .


COPY wait-for.sh .
RUN chmod +x wait-for.sh

EXPOSE 50051
EXPOSE 8080



CMD ["./wait-for.sh", "postgres:5432", "./user-service"]

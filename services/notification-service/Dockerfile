# syntax=docker/dockerfile:1
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY ../../go.mod ../../go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/services/notification-service
RUN go build -o notification-service main.go

FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/services/notification-service/notification-service .

EXPOSE 8086
CMD ["./notification-service"]




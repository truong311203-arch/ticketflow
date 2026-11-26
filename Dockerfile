# Multi-stage Dockerfile for TicketFlow API
# Stage 1: Build the statically linked Go binary
FROM golang:1.27-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/bin/ticketflow-api \
    ./cmd/api

# Stage 2: Minimal, secure runtime image
FROM alpine:3.20

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/bin/ticketflow-api /app/ticketflow-api

USER appuser

EXPOSE 8080

ENTRYPOINT ["/app/ticketflow-api"]

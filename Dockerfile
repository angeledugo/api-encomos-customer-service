# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
RUN apk update && apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN adduser -D -s /bin/sh appuser
WORKDIR /root/
COPY --from=builder /app/main .
RUN chown appuser:appuser main
USER appuser
EXPOSE 50055
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 CMD echo "Health check not implemented yet"
CMD ["./main"]

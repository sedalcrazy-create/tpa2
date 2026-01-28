# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy source code first
COPY . .

# Download dependencies with retry
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=off
RUN for i in 1 2 3 4 5; do go mod tidy && go mod download && break || sleep 5; done

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/tpa-api \
    ./cmd/api

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/tpa-api .

# Set timezone
ENV TZ=Asia/Tehran

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./tpa-api"]

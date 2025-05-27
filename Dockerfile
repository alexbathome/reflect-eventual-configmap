# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files if they exist
COPY go.* ./

# Download dependencies
RUN go mod download || true

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port (if needed)
EXPOSE 8080

# Command to run
ENTRYPOINT ["./main"]
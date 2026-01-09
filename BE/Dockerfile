# Build Stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# Use CGO_ENABLED=0 for static linking suitable for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/server ./cmd/api

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
# Verify we copy it to bin/server to match the CMD path
COPY --from=builder /app/bin/server ./bin/server

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./bin/server"]

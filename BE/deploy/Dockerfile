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
# -o specifies the output path relative to WORKDIR
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/server ./cmd/api

# Run Stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from the builder stage
COPY --from=builder /app/bin/server .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./server"]

# Build stage
FROM golang:1.26-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git and other build dependencies
RUN apk add --update --no-cache git ca-certificates tzdata && update-ca-certificates

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# We build the binary in cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/main.go

# Run stage
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose gRPC port
EXPOSE 50051

# Command to run the executable
ENTRYPOINT ["./main"]

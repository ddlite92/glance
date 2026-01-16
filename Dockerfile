# Use the official Golang image as the build environment
FROM golang:1.24 as builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o glance ./main.go

# Use a minimal image for the final runtime
FROM debian:bookworm-slim

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/glance /app/glance

# Expose the port the app runs on (change if needed)
EXPOSE 8081

# Run the binary
ENTRYPOINT ["/app/glance"]

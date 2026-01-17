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

FROM debian:bookworm-slim

# Install curl for debugging and API testing
RUN apt-get update && apt-get install -y --no-install-recommends curl \
	&& rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/glance /app/glance

# Expose the port the app runs on (change if needed)
EXPOSE 8081

# Run the binary
ENTRYPOINT ["/app/glance"]

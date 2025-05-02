# Use official Go base image
FROM golang:1.24-alpine

# Set environment
ENV GO111MODULE=on

# Create working directory
WORKDIR /app

# Copy Go files
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Build binary
RUN go build -o pwh-auth ./cmd/server

# Expose port (read from env in container)
EXPOSE 5005

# Use entrypoint (can be overridden in docker-compose)
CMD ["./pwh-auth"]

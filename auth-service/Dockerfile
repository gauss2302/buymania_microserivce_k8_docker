# Use the official Golang image for the build stage
FROM golang:1.23-alpine AS builder

# Install required build dependencies
RUN apk add --no-cache gcc musl-dev

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/api/main.go

# Use a minimal image for the final stage
FROM alpine:3.18

# Install required runtime packages
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/auth-service .

# Expose the port the service will run on
EXPOSE 8084

# Start the Auth service
CMD ["./auth-service"]

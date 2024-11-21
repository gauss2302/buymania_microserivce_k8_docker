FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o user-service main.go

# Use a minimal image for the final build
FROM alpine:3.18

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /app/user-service /app/user-service

# Expose the port the service will run on
EXPOSE 8083

# Start the User service
ENTRYPOINT ["/app/user-service"]
# product-service/Dockerfile
FROM golang:1.23-alpine AS builder

# Install required packages
RUN apk add --no-cache gcc musl-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application (путь изменен согласно вашей структуре)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o product-service ./cmd/server/main.go

# Use a minimal image for the final build
FROM alpine:3.18

# Install required runtime packages
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /app/product-service .

# Copy migrations if needed
COPY --from=builder /app/migrations /app/migrations


# Expose the port the service will run on
EXPOSE 8082

# Start the Product service
CMD ["./product-service"]
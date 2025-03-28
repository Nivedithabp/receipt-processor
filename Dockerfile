# Use the correct Go version to build the app
FROM golang:1.24.1 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Install swag CLI tool to generate Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init

# Set environment variables to ensure correct architecture for binary
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Build the Go app for Linux (64-bit)
RUN go build -o /app/main .

# Use a minimal base image for deployment
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary and docs from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

# Add ca-certificates for HTTPS support (optional but recommended)
RUN apk --no-cache add ca-certificates

# Expose the application on port 8080
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

# Use official Golang image
FROM golang:1.21

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy project files
COPY . .

# Build the application
RUN go build -o receipt-processor

# Expose the application on port 8080
EXPOSE 8080

# Run the application
CMD ["./receipt-processor"]

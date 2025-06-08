# Use official Golang image
FROM golang:1.20

# Set working directory
WORKDIR /app

# Copy Go module files (if any)
COPY go.mod ./
COPY go.sum ./
RUN go mod download || true

# Copy the rest of the app
COPY . .

# Build the Go application
RUN go build -o chat-server server.go

# Expose the TCP port
EXPOSE 8080

# Run the server binary
CMD ["./chat-server"]

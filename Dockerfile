# Use the official Golang image to build the application
FROM golang:1.21 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application with static linking
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Use a minimal base image for the final container
FROM debian:bullseye-slim

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Clean up build cache to reduce image size
RUN rm -rf /root/.cache

# Expose the port the application listens on
EXPOSE 7171

# Add labels for versioning and metadata
LABEL version="1.0" \
      description="Go CRUD application" \
      maintainer="your-email@example.com"

# Run the application
CMD ["./main"]

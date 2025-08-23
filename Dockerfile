# Multi-stage Dockerfile for CarZone application
FROM golang:1.24-alpine

# Set maintainer information
LABEL maintainer="Prateek Kumar <prateekkumar72007@gmail.com>"

# Copy Go module files first for better Docker layer caching
COPY go.mod go.sum ./

# Download dependencies using Go modules
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main

# Expose the port that the application listens on
EXPOSE 8080

# Health check to verify the application is running correctly
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Command to run when container starts
CMD ["./main"]

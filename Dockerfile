# Build stage
FROM golang:latest AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN go build -o main .

# Final stage
FROM minio/minio:latest

# Copy the built binary from the builder stage
COPY --from=builder /app/main /usr/local/bin/

# Set environment variables for MinIO
ENV MINIO_ROOT_USER=minioadmin
ENV MINIO_ROOT_PASSWORD=minioadmin

# Expose ports for MinIO and application
EXPOSE 9000 9001

# CMD instruction to start MinIO and the Go application
CMD ["sh", "-c", "minio server ~/minio --console-address :9001 & /usr/local/bin/main"]

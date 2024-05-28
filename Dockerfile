# Stage 1: Build the Golang application
FROM golang:latest AS builder

WORKDIR /app

# Copy the go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy your source directory
COPY src/ ./src/

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./src

# Stage 2: Setup the runtime container
FROM alpine:latest  

# Install ca-certificates and ffmpeg
RUN apk --no-cache add ca-certificates ffmpeg

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Run the application
CMD ["./main"]

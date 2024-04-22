FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .
FROM minio/minio:latest
COPY --from=builder /app/main /usr/local/bin/
ENV MINIO_ROOT_USER=minioadmin
ENV MINIO_ROOT_PASSWORD=minioadmin

EXPOSE 9000 9001
CMD ["minio server ~/minio --console-address :9001 & /usr/local/bin/main"]

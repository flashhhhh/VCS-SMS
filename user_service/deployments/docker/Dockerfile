# Build stage
FROM golang:alpine AS builder

# Install necessary dependencies for building
RUN apk add --no-cache ca-certificates git && \
    update-ca-certificates

WORKDIR /app

# Copy only necessary files for dependency resolution first
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the applications
RUN go build -o rest-server cmd/server/rest/main.go && \
    go build -o grpc-server cmd/server/grpc/main.go

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS connections
RUN apk add --no-cache ca-certificates && \
    update-ca-certificates

WORKDIR /app

# Copy only the built binaries from the builder stage
COPY --from=builder /app/rest-server /app/rest-server
COPY --from=builder /app/grpc-server /app/grpc-server
COPY configs/ ./configs/

# Run both servers
CMD ["sh", "-c", "./rest-server & ./grpc-server"]
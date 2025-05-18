# Build stage
FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o healthcheck_service cmd/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/healthcheck_service .
COPY configs/ ./configs/

CMD ["./healthcheck_service"]
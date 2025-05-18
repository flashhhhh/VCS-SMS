FROM golang:alpine AS builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o mail_service cmd/server/main.go

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates && update-ca-certificates

WORKDIR /app

COPY --from=builder /app/mail_service .
COPY configs/ ./configs/

CMD ["./mail_service"]
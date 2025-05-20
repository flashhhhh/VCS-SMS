# Build stage
FROM golang:alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o rest-server ./cmd/server/rest && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o grpc-server ./cmd/server/grpc && \
    CGO_ENABLED=0 go build -ldflags="-s -w" -o kafka-server ./cmd/server/kafka

# Runtime stage
FROM alpine

RUN apk add --no-cache tini

WORKDIR /app

COPY --from=builder /app/rest-server .
COPY --from=builder /app/grpc-server .
COPY --from=builder /app/kafka-server .
COPY configs/ ./configs/

ENTRYPOINT ["/sbin/tini", "--"]
CMD ["sh", "-c", "./rest-server & ./grpc-server & ./kafka-server"]

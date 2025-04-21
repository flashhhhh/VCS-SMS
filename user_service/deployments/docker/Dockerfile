FROM golang:alpine

# Install necessary dependencies
RUN apk add --no-cache ca-certificates git && \
    update-ca-certificates

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["sh", "-c", "go run cmd/server/rest/main.go & go run cmd/server/grpc/main.go"]
FROM golang:alpine

RUN apk add --no-cache ca-certificates git && \
    update-ca-certificates

WORKDIR /app

COPY . .

RUN go mod tidy

CMD ["go", "run", "cmd/main.go"]
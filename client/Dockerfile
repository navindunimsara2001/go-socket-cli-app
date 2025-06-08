# Dockerfile.client
FROM golang:1.20

WORKDIR /app

COPY . .

RUN go build -o chat-client client.go

CMD ["./chat-client"]

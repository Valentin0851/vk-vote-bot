FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /vote-bot ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /vote-bot /app/vote-bot
COPY config/config.yaml /app/config/config.yaml

EXPOSE 8080
CMD ["/app/vote-bot"]
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/bot .
COPY --from=builder /app/configs ./configs

CMD ["./bot"]
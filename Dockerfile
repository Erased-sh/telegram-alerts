# Этап 1: Сборка
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

CMD ["./main"]

EXPOSE 8080

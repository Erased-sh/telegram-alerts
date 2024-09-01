FROM docker.io/golang:1.22-alpine as build
WORKDIR /app

COPY ./go.mod ./go.sum /app/
COPY . /app

RUN go build -v -o bin main.go

FROM docker.io/alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=build /app/bin /app/bin

COPY selfsigned.crt /etc/ssl/certs/selfsigned.crt
COPY selfsigned.key /etc/ssl/private/selfsigned.key

ENV TELEGRAM_BOT_TOKEN;

CMD ["/app/bin"]

EXPOSE 8080
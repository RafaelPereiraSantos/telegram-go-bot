FROM golang:1.13-alpine as builder
LABEL stage=intermediate

COPY ./ /app
WORKDIR /app

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telegram_bot cmd/telegram_go_bot/main.go

FROM alpine

WORKDIR /app
COPY --from=builder app/telegram_bot ./app/telegram_bot

EXPOSE 8080

CMD ["./app/telegram_bot"]
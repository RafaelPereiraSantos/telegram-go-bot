# telegram-go-bot
A telegram bot written in go

## How to run locally

### Requirements:
First, you must had created a Telegram bot, to do so follow the official documentation about how to do so https://core.telegram.org/bots.

### Option 1 - Run it in your own machine
- After created a bot on telegram, create a file called ```.env``` with the same keys as the file ```.env.local``` and fill the ```TELEGRAM_SECRET_TOKEN``` variable with the telegram key of your bot.
- Run:
```
go run tidy
go mod vendor
go run cmd/telegram_go_bot/main.go
```

### Option 2 - Run it in a docker container
- With the telegram key of your bot on hands, run the following build docker command passing the key as an environment variable to the docker run command:
```
docker build -t my-telegram-bot .
docker run -e TELEGRAM_SECRET_TOKEN="${my-secret-key}" -e PORT=8080 my-telegram-bot
```
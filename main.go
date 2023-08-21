// main.go

package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"telegram-redminenotifybot/bot"
	"telegram-redminenotifybot/redmine"
)

var RedmineClient *redmine.RedmineClient

func main() {
	bot.InitEnv()

	RedmineClient = redmine.NewRedmineClient(os.Getenv("REDMINE_URL"), os.Getenv("REDMINE_API"))

	err := bot.NewBot(os.Getenv("API_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize bot: %s", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.RedmineBot.API.GetUpdatesChan(u)

	for update := range updates {

		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				bot.SendProjectsList(update.Message.Chat.ID, RedmineClient)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Эта команда мне не известна")
				bot.RedmineBot.API.Send(msg)
			}

		} else if update.CallbackQuery != nil {
			bot.HandleCallbackQuery(update.CallbackQuery, RedmineClient)
		}
	}
}

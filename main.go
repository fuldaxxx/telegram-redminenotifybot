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

	updates := bot.RedmineBot.API.GetUpdatesChan(tgbotapi.NewUpdate(0))

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		if update.Message.Text == "/start" {
			bot.SendProjectsList(chatID, RedmineClient)
		}

		if update.CallbackQuery != nil {
			bot.HandleCallbackQuery(update.CallbackQuery, RedmineClient)
		}
	}
}

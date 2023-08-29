package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"telegram-redminenotifybot/models"
)

func InitEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

var RedmineBot *models.Bot

func NewBot(token string) error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("API_TOKEN"))
	if err != nil {
		return err
	}

	RedmineBot = &models.Bot{
		API: bot,
	}

	bot.Debug = true

	log.Printf("Авторизация в %s", bot.Self.UserName)

	return nil
}

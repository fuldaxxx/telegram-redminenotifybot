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

//func GetComment(client *redmine.RedmineClient, taskID int, chatID int64) {
//	comment, err := client.GetTaskJournals(taskID)
//	if err != nil {
//		return
//	}
//	messageText := "лох"
//
//	msg := tgbotapi.NewMessage(chatID, messageText)
//	_, err = RedmineBot.API.Send(msg)
//	if err != nil {
//		log.Printf("Не удалось отправить сообщение %s", err)
//	}
//}

//func GetNewComment(client *redmine.RedmineClient, projectID string, chatID int64, user database.User) {
//	for {
//		tasks, err := client.GetIssuesForProject()
//	}
//}

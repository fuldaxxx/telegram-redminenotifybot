package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"telegram-redminenotifybot/bot"
	"telegram-redminenotifybot/database"
	"telegram-redminenotifybot/redmine"
)

var RedmineClient *redmine.RedmineClient

func main() {
	bot.InitEnv()
	db, eror := gorm.Open(sqlite.Open("redmine_bot.db"), &gorm.Config{})
	if eror != nil {
		log.Fatalf("Failed to connect to database: %s", eror)
	}

	db.AutoMigrate(&database.User{})
	db.AutoMigrate((&database.Project{}))

	err := bot.NewBot(os.Getenv("API_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize bot: %s", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.RedmineBot.API.GetUpdatesChan(u)

	for update := range updates {
		var user database.User
		var chatId int64
		if update.Message == nil {
			chatId = update.CallbackQuery.Message.Chat.ID
		} else {
			chatId = update.Message.Chat.ID
		}
		db.Where("chat_id = ?", chatId).First(&user)
		RedmineClient = redmine.NewRedmineClient(user.RedmineURL, user.APIKey)

		if update.Message != nil {
			switch update.Message.Command() {
			case "tasks":
				bot.SendProjectsList(update.Message.Chat.ID, RedmineClient)
			case "start":
				chatID := update.Message.Chat.ID
				msg := tgbotapi.NewMessage(chatID, "Введите Redmine URL:")
				bot.RedmineBot.API.Send(msg)

				urlUpdate := <-updates
				redmineURL := urlUpdate.Message.Text

				msg = tgbotapi.NewMessage(chatID, "Введите API Key:")
				bot.RedmineBot.API.Send(msg)

				apiKeyUpdate := <-updates
				apiKey := apiKeyUpdate.Message.Text

				user := database.User{
					ChatID:     chatID,
					RedmineURL: redmineURL,
					APIKey:     apiKey,
				}
				db.Create(&user)

				msg = tgbotapi.NewMessage(chatID, "Теперь введите команду /tasks.")
				bot.RedmineBot.API.Send(msg)
			default:
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Эта команда мне не известна")
				bot.RedmineBot.API.Send(msg)
			}

		} else if update.CallbackQuery != nil {
			projectID := update.CallbackQuery.Data
			err = database.SaveUserProject(db, update.CallbackQuery.Message.Chat.ID, projectID)
			if err != nil {
				log.Printf("Error saving user project: %s", err)
				continue
			}
			bot.HandleCallbackQuery(update.CallbackQuery, RedmineClient, user)
		}
	}
}

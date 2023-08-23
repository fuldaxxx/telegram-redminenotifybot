package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"telegram-redminenotifybot/bot"
	"telegram-redminenotifybot/database"
	"telegram-redminenotifybot/redmine"
	"time"
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
	var user database.User
	var projets database.Project

	bot.StartTaskListeners(db)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.RedmineBot.API.GetUpdatesChan(u)

	for update := range updates {
		var chatId int64
		if update.Message == nil {
			chatId = update.CallbackQuery.Message.Chat.ID
		} else {
			chatId = update.Message.Chat.ID
		}
		db.Where("chat_id = ?", chatId).First(&user)
		db.Where("user_id = ?", chatId).First(&projets)
		RedmineClient = redmine.NewRedmineClient(user.RedmineURL, user.APIKey)

		if update.Message != nil {
			switch update.Message.Command() {
			case "tasks":
				db.Where("user_id = ?", chatId).First(&projets)
				bot.SendTaskList(chatId, RedmineClient, projets.ProjectID, user)
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
			case "projects":
				bot.SendProjectsList(update.Message.Chat.ID, RedmineClient)
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
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID,
				fmt.Sprintf("Теперь вам будут приходить уведомления по проекту: %s", update.CallbackQuery.Data))
			bot.RedmineBot.API.Send(msg)
		}
	}
}

func RunUserTaskChecker(db *gorm.DB, user database.User) {
	RedmineClient := redmine.NewRedmineClient(user.RedmineURL, user.APIKey)

	for {
		var project database.Project
		if err := db.Where("user_id = ?", user.ChatID).First(&project).Error; err != nil {
			log.Printf("Error fetching project record for user %d: %s", user.ChatID, err)
			return
		}

		bot.GetNewTask(RedmineClient, project.ProjectID, user.ChatID, user)

		time.Sleep(time.Minute) // Интервал между выполнениями для данного пользователя
	}
}

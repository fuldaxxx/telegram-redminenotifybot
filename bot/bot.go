package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"telegram-redminenotifybot/models"
	"telegram-redminenotifybot/redmine"
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

func SendProjectsList(chatID int64, RedmineClient *redmine.RedmineClient) {
	projects, err := RedmineClient.GetProjects()
	if err != nil {
		log.Printf("Error fetching projects: %s", err)
		return
	}

	messageText := "Выберите проект:\n"

	for _, project := range projects {
		messageText += strconv.Itoa(project.ID) + " - " + project.Name + "\n"
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	RedmineBot.API.Send(msg)
}

func HandleCallbackQuery(query *tgbotapi.CallbackQuery, RedmineClient *redmine.RedmineClient) {
	projectID, err := strconv.Atoi(query.Data)
	if err != nil {
		log.Printf("Error parsing project ID: %s", err)
		return
	}

	issues, err := RedmineClient.GetIssuesForProject(projectID)
	if err != nil {
		log.Printf("Error fetching issues: %s", err)
		return
	}

	messageText := "Список задач:\n"

	for _, issue := range issues {
		messageText += "Номер: " + strconv.Itoa(issue.ID) + "\n"
		messageText += "Тема: " + issue.Subject + "\n"
		messageText += "Описание: " + issue.Description + "\n\n"
	}

	msg := tgbotapi.NewMessage(query.Message.Chat.ID, messageText)
	RedmineBot.API.Send(msg)
}

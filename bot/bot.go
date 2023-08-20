package bot

import (
	"fmt"
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
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, project := range projects {
		messageText += strconv.Itoa(project.ID) + " - " + project.Name + "\n"

		btn := tgbotapi.NewInlineKeyboardButtonData(project.Name, strconv.Itoa(project.ID))
		row := tgbotapi.NewInlineKeyboardRow(btn)
		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = keyboard
	RedmineBot.API.Send(msg)
}

func SendTaskList(chatID int64, RedmineClient *redmine.RedmineClient, projectID int) {
	tasks, err := RedmineClient.GetIssuesForProject(projectID)
	if err != nil {
		log.Printf("Не удалось отправить задачи по проекту %d: %s", projectID, err)
	}

	messageText := fmt.Sprintf("Задачи по проекту %d\n", projectID)

	for _, task := range tasks {
		messageText += fmt.Sprintf("Номер задачи: %d\n Автор задачи: %s\n Статус задачи: %s\n Тема: %s\n, Описание: %s",
			task.ID, task.Author, task.Status, task.Subject, task.Description)
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

	SendTaskList(query.Message.Chat.ID, RedmineClient, projectID)
}

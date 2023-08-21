package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"golang.org/x/net/html"
	"log"
	"os"
	"strconv"
	"strings"
	"telegram-redminenotifybot/database"
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

func cleanHTMLTags(htmlText string) string {
	tokenizer := html.NewTokenizer(strings.NewReader(htmlText))
	cleanText := ""

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return cleanText
		case html.TextToken:
			cleanText += tokenizer.Token().Data
		}
	}
}

func SendProjectsList(chatID int64, RedmineClient *redmine.RedmineClient) {
	projects, err := RedmineClient.GetProjects()
	if err != nil {
		log.Printf("Error fetching projects: %s", err)
		return
	}

	messageText := "Выберите проект\n"
	var rows [][]tgbotapi.InlineKeyboardButton

	for _, project := range projects {
		btn := tgbotapi.NewInlineKeyboardButtonData(project.Name, strconv.Itoa(project.ID))
		row := tgbotapi.NewInlineKeyboardRow(btn)
		rows = append(rows, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ReplyMarkup = keyboard
	RedmineBot.API.Send(msg)
}

func SendTaskList(chatID int64, RedmineClient *redmine.RedmineClient, projectID int, user database.User) {
	redmineURL := user.RedmineURL
	tasks, err := RedmineClient.GetIssuesForProject(projectID)
	if err != nil {
		log.Printf("Не удалось получить задачи по проекту %d: %s", projectID, err)
		errorMsg := fmt.Sprintf("Произошла ошибка при получении задач для проекта %d", projectID)
		msg := tgbotapi.NewMessage(chatID, errorMsg)
		RedmineBot.API.Send(msg)
		return
	}

	messageText := fmt.Sprintf("Задачи по проекту %d\n", projectID)

	for _, task := range tasks {
		if task.Status.Name != "Решена" && task.Status.Name != "Обратная связь" && len(task.Description) > 5 {
			taskText := fmt.Sprintf("Номер задачи: %d\nАвтор задачи: %s\nСтатус задачи: %s\nТема: %s\n",
				task.ID, task.Author.Name, task.Status.Name, task.Subject)

			messageText += taskText + "\n"
			messageText += fmt.Sprintf("%s/issues/%d", redmineURL, task.ID)
			messageText += "\n\n"
		}
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	_, err = RedmineBot.API.Send(msg)
	if err != nil {
		log.Printf("Error sending message: %s", err)
	}

}

func HandleCallbackQuery(query *tgbotapi.CallbackQuery, RedmineClient *redmine.RedmineClient, user database.User) {
	projectID, err := strconv.Atoi(query.Data)
	if err != nil {
		log.Printf("Error parsing project ID: %s", err)
		return
	}

	log.Printf(query.Data)

	SendTaskList(query.Message.Chat.ID, RedmineClient, projectID, user)
}

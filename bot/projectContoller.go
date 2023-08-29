package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"telegram-redminenotifybot/models"
)

func SendProjectsList(chatID int64, RedmineClient *models.RedmineClient) {
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

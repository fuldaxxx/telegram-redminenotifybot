package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"telegram-redminenotifybot/database"
	"telegram-redminenotifybot/models"
	"time"
)

func GetNewComments(RedmineClient *models.RedmineClient, chatID int64, user database.User, projectID string) {
	interval := time.Minute

	lastComment := make(map[int][]models.Journals)

	for {
		tasks, err := RedmineClient.GetIssuesForProject(projectID)
		if err != nil {
			fmt.Printf("Ошибка получения задач: %s", err)
		}

		for _, task := range tasks {
			comments, err := RedmineClient.GetTaskJournals(task.ID)
			if err != nil {
				fmt.Printf("Ошибка получения комментариев: %s", err)
			}

			lastCommentForTask := lastComment[task.ID]
			newComments := findNewComments(lastCommentForTask, comments)

			messageText := fmt.Sprintf("⚡️✉️ Новый комментарий к задаче\n\n")
			if len(newComments) != 0 {
				taskText := fmt.Sprintf("Номер: %d\nАвтор задачи: %s\nСтатус задачи: %s\nТема: %s\n\n",
					task.ID, task.Author.Name, task.Status.Name, task.Subject)

				messageText += taskText
				messageText += fmt.Sprintf("%s/issues/%d#note-%d", user.RedmineURL, task.ID, len(comments))
				messageText += "\n\n"

			}

			if messageText != fmt.Sprintf("⚡️✉️ Новый комментарий к задаче\n\n") && len(lastCommentForTask) != 0 {
				msg := tgbotapi.NewMessage(chatID, messageText)
				_, err = RedmineBot.API.Send(msg)
				if err != nil {
					log.Printf("Error sending message: %s", err)
				}
			}

			lastComment[task.ID] = comments
		}

		time.Sleep(interval)

	}
}

func findNewComments(lastComments []models.Journals, newComments []models.Journals) []models.Journals {
	var commentsArray []models.Journals

	for _, newComment := range newComments {
		found := false
		for _, lastComment := range lastComments {
			if newComment.ID == lastComment.ID {
				found = true
			}
		}

		if !found {
			commentsArray = append(commentsArray, newComment)
		}
	}

	return commentsArray
}

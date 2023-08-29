package bot

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"log"
	"telegram-redminenotifybot/database"
	"telegram-redminenotifybot/models"
	"time"
)

func SendTaskList(chatID int64, RedmineClient *models.RedmineClient, projectID string, user database.User) {
	redmineURL := user.RedmineURL
	tasks, err := RedmineClient.GetIssuesForProject(projectID)
	if err != nil {
		log.Printf("Не удалось получить задачи по проекту %s: %s", projectID, err)
		errorMsg := fmt.Sprintf("Произошла ошибка при получении задач для проекта %s", projectID)
		msg := tgbotapi.NewMessage(chatID, errorMsg)
		RedmineBot.API.Send(msg)
		return
	}

	messageText := fmt.Sprintf("Задачи по проекту %s\n", projectID)

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

func GetNewTask(client *models.RedmineClient, projectID string, chatID int64, user database.User) {
	interval := time.Minute

	var lastTask []models.Issue

	for {
		tasks, err := client.GetIssuesForProject(projectID)
		if err != nil {
			fmt.Printf("Ошибка получения обновления о задачах: %s", err)
		}

		newIssues := findNewIssue(lastTask, tasks)
		messageText := fmt.Sprintf("Новая задача по проекту: %s\n", projectID)
		for _, newIssue := range newIssues {
			if newIssue.Status.Name != "Решена" && len(newIssue.Description) > 5 {
				taskText := fmt.Sprintf("Номер задачи: %d\nАвтор задачи: %s\nСтатус задачи: %s\nТема: %s\n",
					newIssue.ID, newIssue.Author.Name, newIssue.Status.Name, newIssue.Subject)

				messageText += taskText + "\n"
				messageText += fmt.Sprintf("%s/issues/%d", user.RedmineURL, newIssue.ID)
				messageText += "\n\n"
			}
		}

		if messageText != fmt.Sprintf("Новая задача по проекту: %s\n", projectID) && len(lastTask) != 0 {
			msg := tgbotapi.NewMessage(chatID, messageText)
			_, err = RedmineBot.API.Send(msg)
			if err != nil {
				log.Printf("Error sending message: %s", err)
			}
		}

		lastTask = tasks

		time.Sleep(interval)
	}

}

func findNewIssue(lastTasks []models.Issue, newTasks []models.Issue) []models.Issue {
	var tasksArray []models.Issue

	for _, newTask := range newTasks {
		found := false
		for _, lastTask := range lastTasks {
			if newTask.ID == lastTask.ID {
				found = true
			}
		}

		if !found {
			tasksArray = append(tasksArray, newTask)
		}
	}

	return tasksArray
}

func StartTaskListeners(db *gorm.DB) {
	var projects []database.Project
	db.Find(&projects)

	for _, project := range projects {
		go func(p database.Project) {
			user := getUserByChatID(db, p.UserID)
			if user == nil {
				log.Printf("Пользователь не найден для проекта: %+v", p)
				return
			}

			RedmineClient := NewRedmineClient(user.RedmineURL, user.APIKey)
			GetNewTask(RedmineClient, p.ProjectID, user.ChatID, *user)
		}(project)
	}
}

func getUserByChatID(db *gorm.DB, chatID int64) *database.User {
	user := &database.User{}
	result := db.Unscoped().Where("chat_id = ?", chatID).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Printf("Ошибка при получении записи о пользователе: %s", result.Error)
		return nil
	}
	return user
}

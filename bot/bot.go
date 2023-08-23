package bot

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"telegram-redminenotifybot/database"
	"telegram-redminenotifybot/models"
	"telegram-redminenotifybot/redmine"
	"time"
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

func SendTaskList(chatID int64, RedmineClient *redmine.RedmineClient, projectID string, user database.User) {
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

func HandleCallbackQuery(query *tgbotapi.CallbackQuery, RedmineClient *redmine.RedmineClient, user database.User) {
	projectID := query.Data
	SendTaskList(query.Message.Chat.ID, RedmineClient, projectID, user)
}

func GetNewTask(client *redmine.RedmineClient, projectID string, chatID int64, user database.User) {
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
			if newIssue.Status.Name != "Решена" && newIssue.Status.Name != "Обратная связь" && len(newIssue.Description) > 5 {
				taskText := fmt.Sprintf("Номер задачи: %d\nАвтор задачи: %s\nСтатус задачи: %s\nТема: %s\n",
					newIssue.ID, newIssue.Author.Name, newIssue.Status.Name, newIssue.Subject)

				messageText += taskText + "\n"
				messageText += fmt.Sprintf("%s/issues/%d", user.RedmineURL, newIssue.ID)
				messageText += "\n\n"
			}
		}

		if messageText != fmt.Sprintf("Новая задача по проекту: %s\n", projectID) {
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

			RedmineClient := redmine.NewRedmineClient(user.RedmineURL, user.APIKey)
			GetNewTask(RedmineClient, p.ProjectID, user.ChatID, *user)
		}(project)
	}
}

func getUserByChatID(db *gorm.DB, chatID int64) *database.User {
	user := &database.User{}
	result := db.Where("chat_id = ?", chatID).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Printf("Ошибка при получении записи о пользователе: %s", result.Error)
		return nil
	}
	return user
}

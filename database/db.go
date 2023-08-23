package database

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ChatID     int64 `gorm:"uniqueIndex"`
	RedmineURL string
	APIKey     string
	Project    Project
}

type Project struct {
	gorm.Model
	ProjectID string
	UserID    int64
}

var DB *gorm.DB

func SaveUserProject(db *gorm.DB, chatID int64, projectID string) error {
	fmt.Printf("chatID: %d, projectID: %s\n", chatID, projectID)

	user := &User{}
	err := db.Where("chat_id = ?", chatID).First(user).Error
	if err != nil {
		fmt.Printf("Error fetching user record: %s\n", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user = &User{ChatID: chatID, Project: Project{ProjectID: projectID}}
			return db.Create(user).Error
		}
		return err
	}

	var project Project
	err = db.Where("user_id = ?", chatID).First(&project).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// Если не существует, создать новую запись
		project = Project{UserID: chatID, ProjectID: projectID}
		return db.Create(&project).Error
	} else if err != nil {
		// Если произошла другая ошибка, вернуть ее
		fmt.Printf("Error checking project record: %s\n", err)
		return err
	}

	// Выполните обновление проекта в таблице projects
	err = db.Model(&Project{}).Where("user_id = ?", chatID).
		Update("project_id", projectID).Error
	if err != nil {
		fmt.Printf("Error updating project record: %s\n", err)
		return err
	}

	return nil
}

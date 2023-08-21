package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
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
	user := User{ChatID: chatID}
	project := Project{ProjectID: projectID, UserID: chatID}

	err := db.FirstOrCreate(&user, user).Error
	if err != nil {
		return err
	}

	return db.FirstOrCreate(&project, project).Error
}

func OpenDatabase() (*gorm.DB, error) {
	dbFile := "redmine_bot.db"

	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		return gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	}

	return gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
}

func InitDB() error {
	db, err := gorm.Open(sqlite.Open("redmine_bot.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	DB = db

	db.AutoMigrate(&User{})

	return nil
}

func GetUserByUsername(username string) (User, error) {
	var user User
	result := DB.Where("telegram_username = ?", username).First(&user)
	if result.Error != nil {
		return User{}, result.Error
	}

	return user, nil
}

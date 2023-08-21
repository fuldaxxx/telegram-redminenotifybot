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
}

var DB *gorm.DB

func OpenDatabase() (*gorm.DB, error) {
	dbFile := "redmine_bot.db"

	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	if os.IsNotExist(err) {
		// Файл базы данных не существует, создаем новую базу данных
		return gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	}

	// Файл базы данных существует, открываем его
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

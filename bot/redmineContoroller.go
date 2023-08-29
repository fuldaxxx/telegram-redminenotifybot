package bot

import (
	"net/http"
	"telegram-redminenotifybot/models"
)

func NewRedmineClient(url, token string) *models.RedmineClient {
	return &models.RedmineClient{
		URL:    url,
		Token:  token,
		Client: &http.Client{},
	}
}

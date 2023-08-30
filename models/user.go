package models

import "time"

type InfoAboutUser struct {
	User UserAccount `json:"user"`
}

type UserAccount struct {
	ID          int       `json:"id"`
	Login       string    `json:"login"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	Mail        string    `json:"mail"`
	CreatedOn   time.Time `json:"created_on"`
	LastLoginOn time.Time `json:"last_login_on"`
}

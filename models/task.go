package models

import "fmt"

type IssueList struct {
	Issue []Issue `json:"issues"`
}

type Issue struct {
	ID          int        `json:"id"`
	Project     Projects   `json:"project"`
	Tracker     Tracker    `json:"tracker"`
	Status      Status     `json:"status"`
	Priority    Priority   `json:"priority"`
	Author      Author     `json:"author"`
	AssignedTo  AssignedTo `json:"assigned_to"`
	Subject     string     `json:"subject"`
	Description string     `json:"description"`
}

type Tracker struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Priority struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AssignedTo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (i Issue) GetTask() string {
	return fmt.Sprintf("Задача #%d \nТема: %s \nОписание: %s \n Статус: %s",
		i.ID, i.Subject, i.Description, i.Status)
}

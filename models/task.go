package models

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
	Journals    []Journals `json:"journals"`
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

type Journals struct {
	ID        int    `json:"id"`
	User      User   `json:"user"`
	Notes     string `json:"notes"`
	CreatedOn string `json:"created_on"`
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

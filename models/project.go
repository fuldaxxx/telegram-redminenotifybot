package models

type ProjectsList struct {
	Projects []Projects `json:"projects"`
}

type Projects struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

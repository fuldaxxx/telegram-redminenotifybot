package models

import "fmt"

type ProjectsList struct {
	Projects []Projects `json:"projects"`
}

type Projects struct {
	ID   int `json:"id"`
	Name int `json:"name"`
}

func (p Projects) GetProjects() string {
	return fmt.Sprintf("%s - %d", p.Name, p.ID)
}

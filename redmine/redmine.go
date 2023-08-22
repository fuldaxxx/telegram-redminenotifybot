package redmine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"telegram-redminenotifybot/models"
)

type RedmineClient struct {
	URL    string
	Token  string
	Client *http.Client
}

func NewRedmineClient(url, token string) *RedmineClient {
	return &RedmineClient{
		URL:    url,
		Token:  token,
		Client: &http.Client{},
	}
}

func (r *RedmineClient) GetIssuesForProject(projectID string) ([]models.Issue, error) {
	url := fmt.Sprintf("%s/issues.json?project_id=%s", r.URL, projectID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Redmine-API-Key", r.Token)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var issues models.IssueList
	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return nil, err
	}

	return issues.Issue, nil
}

func (r *RedmineClient) GetProjects() ([]models.Projects, error) {
	url := fmt.Sprintf("%s/projects.json", r.URL)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Redmine-API-Key", r.Token)

	resp, err := r.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projectsList models.ProjectsList
	err = json.NewDecoder(resp.Body).Decode(&projectsList)
	if err != nil {
		return nil, err
	}

	return projectsList.Projects, nil
}

package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RedmineClient struct {
	URL    string
	Token  string
	Client *http.Client
}

func (r *RedmineClient) GetIssuesForProject(projectID string) ([]Issue, error) {
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

	var issues IssueList
	err = json.NewDecoder(resp.Body).Decode(&issues)
	if err != nil {
		return nil, err
	}

	return issues.Issue, nil
}

func (r *RedmineClient) GetProjects() ([]Projects, error) {
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

	var projectsList ProjectsList
	err = json.NewDecoder(resp.Body).Decode(&projectsList)
	if err != nil {
		return nil, err
	}

	return projectsList.Projects, nil
}

func (r *RedmineClient) GetTaskJournals(taskID int) ([]Journals, error) {
	url := fmt.Sprintf("/issues/%d.json?incldue=journals", taskID)

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

	var issue Issue
	err = json.NewDecoder(resp.Body).Decode(&issue)
	if err != nil {
		return nil, err
	}

	return issue.Journals, nil
}

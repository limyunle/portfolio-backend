package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var httpGet = http.Get

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"html_url"`
}

type Service interface {
	GetRepos(username string) ([]Repo, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) GetRepos(username string) ([]Repo, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)

	resp, err := httpGet(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var repos []Repo
	if err := json.Unmarshal(body, &repos); err != nil {
		return nil, fmt.Errorf("failed to parse repos: %w", err)
	}
	return repos, nil
}

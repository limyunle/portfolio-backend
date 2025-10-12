package leetcode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var httpGet = http.Get

type LeetCodeStats struct {
	TotalSolved        int            `json:"totalSolved"`
	EasySolved         int            `json:"easySolved"`
	MediumSolved       int            `json:"mediumSolved"`
	HardSolved         int            `json:"hardSolved"`
	SubmissionCalendar map[string]int `json:"submissionCalendar"`
}

type Service interface {
	GetLeetCodeStats(username string) (*LeetCodeStats, error)
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) GetLeetCodeStats(username string) (*LeetCodeStats, error) {
	url := fmt.Sprintf("https://leetcode-stats-api.herokuapp.com/%s", username)
	resp, err := httpGet(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch leetcode stats: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var stats LeetCodeStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse leetcode stats: %w", err)
	}

	return &stats, nil
}

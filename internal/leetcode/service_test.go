package leetcode

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLeetCodeStats_Success(t *testing.T) {
	mockResponse := LeetCodeStats{
		TotalSolved:  300,
		EasySolved:   100,
		MediumSolved: 150,
		HardSolved:   50,
		SubmissionCalendar: map[string]int{
			"1697068800": 5,
			"1697155200": 7,
		},
	}
	responseBody, _ := json.Marshal(mockResponse)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(responseBody)
	}))
	defer ts.Close()

	s := &service{}
	username := "testuser"

	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL + "/" + username)
	}
	defer func() { httpGet = http.Get }()

	stats, err := s.GetLeetCodeStats(username)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.TotalSolved != 300 {
		t.Errorf("expected TotalSolved=300, got %d", stats.TotalSolved)
	}
	if stats.SubmissionCalendar["1697068800"] != 5 {
		t.Errorf("expected submission 5, got %d", stats.SubmissionCalendar["1697068800"])
	}
}

func TestGetLeetCodeStats_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "{invalid json}")
	}))
	defer ts.Close()

	s := &service{}
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(ts.URL + "/testuser")
	}
	defer func() { httpGet = http.Get }()

	_, err := s.GetLeetCodeStats("testuser")
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestGetLeetCodeStats_HTTPError(t *testing.T) {
	s := &service{}

	httpGet = func(url string) (*http.Response, error) {
		return nil, fmt.Errorf("network unreachable")
	}
	defer func() { httpGet = http.Get }()

	_, err := s.GetLeetCodeStats("testuser")
	if err == nil {
		t.Fatal("expected error due to HTTP failure")
	}
}

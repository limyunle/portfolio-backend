package github

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRepos(t *testing.T) {
	mockResponse := `[{"name":"repo1","html_url":"https://github.com/limyunle/repo1"},
	                  {"name":"repo2","html_url":"https://github.com/limyunle/repo2"}]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	originalGet := httpGet
	httpGet = func(url string) (*http.Response, error) {
		return http.Get(server.URL)
	}
	defer func() { httpGet = originalGet }()

	s := NewService()
	repos, err := s.GetRepos("limyunle")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(repos))
	}
}

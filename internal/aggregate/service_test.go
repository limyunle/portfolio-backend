package aggregate

import (
	"context"
	"testing"
	"time"

	"github.com/limyunle/portfolio-backend/internal/github"
	"github.com/limyunle/portfolio-backend/internal/leetcode"
)

type mockGitHubService struct{}

func (m *mockGitHubService) GetRepos(username string) ([]github.Repo, error) {
	return []github.Repo{
		{Name: "repo1", URL: "https://github.com/user/repo1"},
		{Name: "repo2", URL: "https://github.com/user/repo2"},
	}, nil
}

type mockLeetCodeService struct{}

func (m *mockLeetCodeService) GetLeetCodeStats(username string) (*leetcode.LeetCodeStats, error) {
	return &leetcode.LeetCodeStats{
		TotalSolved:  10,
		EasySolved:   5,
		MediumSolved: 4,
		HardSolved:   1,
	}, nil
}

type mockS3Service struct {
	uploaded bool
	lastKey  string
	lastData interface{}
}

func (m *mockS3Service) UploadJSON(ctx context.Context, bucket, key string, data interface{}) error {
	m.uploaded = true
	m.lastKey = key
	m.lastData = data
	return nil
}

func (m *mockS3Service) GetJSON(ctx context.Context, bucket, key string, data interface{}) error {
	return nil
}

func TestAggregateService_RefreshAndStore(t *testing.T) {
	mockS3 := &mockS3Service{}

	service := &Service{
		GitHubService:   &mockGitHubService{},
		LeetCodeService: &mockLeetCodeService{},
		S3Service:       mockS3,
		BucketName:      "test-bucket",
	}

	err := service.RefreshAndStore("testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mockS3.uploaded {
		t.Fatalf("expected S3 UploadJSON to be called")
	}

	agg, ok := mockS3.lastData.(CombinedStats)
	if !ok {
		t.Fatalf("expected uploaded data to be of type Stats")
	}

	if len(agg.GitHubRepos) != 2 {
		t.Errorf("expected 2 GitHub repos, got %d", len(agg.GitHubRepos))
	}
	if agg.GitHubRepos[0].URL != "https://github.com/user/repo1" {
		t.Errorf("unexpected first repo URL: %s", agg.GitHubRepos[0].URL)
	}

	if agg.LeetCodeStats.TotalSolved != 10 {
		t.Errorf("expected totalSolved 10, got %d", agg.LeetCodeStats.TotalSolved)
	}
	if agg.LeetCodeStats.EasySolved != 5 {
		t.Errorf("expected easySolved 5, got %d", agg.LeetCodeStats.EasySolved)
	}

	// Validate timestamp
	if time.Since(agg.FetchedAt) > 5*time.Second {
		t.Errorf("expected FetchedAt to be recent, got %v", agg.FetchedAt)
	}
}

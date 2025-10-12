package aggregate

import (
	"context"
	"fmt"
	"time"

	"github.com/limyunle/portfolio-backend/internal/github"
	"github.com/limyunle/portfolio-backend/internal/leetcode"
	"github.com/limyunle/portfolio-backend/internal/s3"
)

type Service struct {
	GitHubService   github.Service
	LeetCodeService leetcode.Service
	S3Service       s3.Storage
	BucketName      string
}

type CombinedStats struct {
	GitHubRepos   []github.Repo           `json:"githubRepos"`
	LeetCodeStats *leetcode.LeetCodeStats `json:"leetcodeStats"`
	FetchedAt     time.Time               `json:"fetchedAt"`
}

func (s *Service) RefreshAndStore(username string) error {
	repos, err := s.GitHubService.GetRepos(username)
	if err != nil {
		return fmt.Errorf("github fetch failed: %w", err)
	}

	stats, err := s.LeetCodeService.GetLeetCodeStats(username)
	if err != nil {
		return fmt.Errorf("leetcode fetch failed: %w", err)
	}

	agg := CombinedStats{
		GitHubRepos:   repos,
		LeetCodeStats: stats,
		FetchedAt:     time.Now().UTC(),
	}

	if err := s.S3Service.UploadJSON(context.Background(), s.BucketName, "aggregate-stats.json", agg); err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}

	return nil
}

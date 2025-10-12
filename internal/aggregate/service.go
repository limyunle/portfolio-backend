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
	startTotal := time.Now()
	fmt.Println("Starting RefreshAndStore...")

	start := time.Now()
	githubRepos, err := s.GitHubService.GetRepos(username)
	if err != nil {
		return fmt.Errorf("github fetch failed: %w", err)
	}
	fmt.Printf("GitHub fetch took %s\n", time.Since(start))

	start = time.Now()
	leetcodeStats, err := s.LeetCodeService.GetLeetCodeStats(username)
	if err != nil {
		return fmt.Errorf("leetcode fetch failed: %w", err)
	}
	fmt.Printf("LeetCode fetch took %s\n", time.Since(start))

	combinedStats := CombinedStats{
		GitHubRepos:   githubRepos,
		LeetCodeStats: leetcodeStats,
		FetchedAt:     time.Now().UTC(),
	}

	start = time.Now()
	if err := s.S3Service.UploadJSON(context.Background(), s.BucketName, "aggregate-stats.json", combinedStats); err != nil {
		return fmt.Errorf("failed to upload to S3: %w", err)
	}
	fmt.Printf("S3 upload took %s\n", time.Since(start))

	fmt.Printf("RefreshAndStore completed in total: %s\n", time.Since(startTotal))
	return nil
}

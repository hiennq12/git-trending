package github

import (
	"fmt"
	"github_trending/internal/models"
	"github_trending/internal/telegram"
	"log"
	"time"
)

type TrendingJob struct {
	githubService  *Service
	telegramClient *telegram.Client
}

func NewTrendingJob(githubService *Service, telegramClient *telegram.Client) *TrendingJob {
	return &TrendingJob{
		githubService:  githubService,
		telegramClient: telegramClient,
	}
}

func (j *TrendingJob) Run() error {
	log.Printf("Starting trending job at %v", time.Now().Format("2006-01-02 15:04:05"))

	// Get trending repositories
	repos, err := j.githubService.GetTrending(models.TrendingOptions{})
	if err != nil {
		return fmt.Errorf("failed to get trending repos: %w", err)
	}

	// Build and send message
	message := telegram.BuildMessage(repos)
	if err := j.telegramClient.SendMessage(message); err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	log.Printf("Completed trending job at %v", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

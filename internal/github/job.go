package github

import (
	"fmt"
	"github_trending/internal/models"
	"github_trending/internal/openai"
	"github_trending/internal/telegram"
	"log"
	"time"
)

type TrendingJob struct {
	githubService  *Service
	telegramClient *telegram.Client
	openaiClient   *openai.Client
}

func NewTrendingJob(githubService *Service, telegramClient *telegram.Client, openaiClient *openai.Client) *TrendingJob {
	return &TrendingJob{
		githubService:  githubService,
		telegramClient: telegramClient,
		openaiClient:   openaiClient,
	}
}

func (j *TrendingJob) Run() error {
	log.Printf("Starting trending job at %v", time.Now().Format("2006-01-02 15:04:05"))

	// Get trending repositories
	repos, err := j.githubService.GetTrending(models.TrendingOptions{})
	if err != nil {
		return fmt.Errorf("failed to get trending repos: %w", err)
	}

	// Enhance descriptions with OpenAI
	for _, repo := range repos {
		repoFullName := fmt.Sprintf("%s/%s", repo.Author, repo.Name)
		repoInfo := fmt.Sprintf("Repository: %s\nOriginal Description: %s\nLanguage: %s\nStars: %d",
			repoFullName, repo.Description, repo.Language, repo.Stars)

		enhancedDesc, err := j.openaiClient.GenerateDescription(repoFullName, repoInfo)
		if err != nil {
			log.Printf("Error generating description for %s: %v", repo.Name, err)
			continue
		}

		repo.EnhancedDescription = enhancedDesc
		// Add small delay only for new descriptions (when not from cache)
		if enhancedDesc != "" {
			time.Sleep(time.Second)
		}
	}

	// Build and send message
	// for all repos then build messages and send to telegram one by one
	// since the maximum size of each message sent in telegram is 4096 characters

	for i, repo := range repos {
		message := telegram.BuildMessage(repo, i)
		err = j.telegramClient.SendMessage(message)

		if err != nil {
			return fmt.Errorf("failed to send telegram message: %w", err)
		}
	}

	err = j.telegramClient.SendMessage("========END DAY========")

	if err != nil {
		return fmt.Errorf("failed to send end day message: %w", err)
	}

	log.Printf("Completed trending job at %v", time.Now().Format("2006-01-02 15:04:05"))
	return nil
}

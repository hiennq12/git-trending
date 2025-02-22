// File: cmd/main.go
package main

import (
	"github_trending/internal/cache"
	"github_trending/internal/github"
	"github_trending/internal/openai"
	"github_trending/internal/scheduler"
	"github_trending/internal/telegram"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

func main() {
	// Initialize cache
	dbPath := filepath.Join("data", "cache.db")
	os.MkdirAll("data", 0755) // Create data directory if not exists
	cacheDB, err := cache.NewSQLiteCache(dbPath)
	if err != nil {
		log.Fatalf("Error initializing cache: %v", err)
	}
	defer cacheDB.Close()

	// Initialize services
	githubService := github.NewService()

	telegramConfig := telegram.NewConfig()
	telegramClient, err := telegram.NewClient(telegramConfig)
	if err != nil {
		log.Fatalf("Error initializing Telegram client: %v", err)
	}

	openaiConfig := openai.NewConfig()
	openaiClient := openai.NewClient(openaiConfig, cacheDB)

	// Initialize scheduler
	scheduler := scheduler.NewScheduler()

	// Create trending job
	trendingJob := github.NewTrendingJob(githubService, telegramClient, openaiClient)

	// Schedule job to run at 8:30 AM daily
	// 30 8 * * *
	// */30 * * * *
	//err = scheduler.AddJob("30 8 * * *", trendingJob)
	//if err != nil {
	//	log.Fatalf("Error adding job: %v", err)
	//}

	// Start scheduler
	scheduler.Start()
	log.Println("Scheduler started. Job will run daily at 8:30 AM (ICT)...")

	// Run job immediately for testing
	if err := trendingJob.Run(); err != nil {
		log.Printf("Error running initial job: %v", err)
	}

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down scheduler...")
	scheduler.Stop()
	log.Println("Scheduler stopped successfully")
}

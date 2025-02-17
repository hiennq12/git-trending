package openai

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"github_trending/internal/cache"
	"log"
)

type Client struct {
	client *openai.Client
	cache  *cache.SQLiteCache
}

func NewClient(config *Config, cacheDB *cache.SQLiteCache) *Client {
	return &Client{
		client: openai.NewClient(config.APIKey),
		cache:  cacheDB,
	}
}

func (c *Client) GenerateDescription(repoFullName, repoInfo string) (string, error) {
	// Try to get from cache first
	cachedDesc, err := c.cache.Get(repoFullName)
	if err != nil {
		log.Printf("Error checking cache for %s: %v", repoFullName, err)
	}

	// If found in cache, return it
	if cachedDesc != "" {
		log.Printf("Cache hit for repository: %s", repoFullName)
		return cachedDesc, nil
	}

	// If not in cache, generate new description
	log.Printf("Cache miss for repository: %s, generating new description \n", repoFullName)

	ctx := context.Background()
	prompt := fmt.Sprintf(`Given the following GitHub repository information, 
    provide a detailed description in about 100 words that explains the purpose, 
    features, and potential use cases of the repository:

    %s

    Focus on technical aspects and practical applications.`, repoInfo)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 200,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error generating description: %w", err)
	}

	description := resp.Choices[0].Message.Content
	// Save to cache
	err = c.cache.Set(repoFullName, description)
	if err != nil {
		log.Printf("Error saving to cache for %s: %v", repoFullName, err)
	}

	return description, nil
}

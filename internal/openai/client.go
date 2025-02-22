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
	cachedDesc, err := c.cache.GetDescription(repoFullName)
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
    provide a detailed description that explains the repository in the following format:

    %s

    Please structure the response as follows:
    1. Overview (1-2 sentences about what the repository is)
    2. Key Features:
       - List 3-4 main features with technical details
       - Each bullet point should be 1-2 sentences
    3. Use Cases: 
       - List 2-3 practical applications
       - Each bullet point should be 1 sentence

    Focus on technical aspects and practical applications. Keep the total response within 100 words.`, repoInfo)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: 300,
		},
	)

	if err != nil {
		return "", fmt.Errorf("error generating description: %w", err)
	}

	description := resp.Choices[0].Message.Content
	// Save to cache
	err = c.cache.SetDescripton(repoFullName, description)
	if err != nil {
		log.Printf("Error saving to cache for %s: %v", repoFullName, err)
	}

	return description, nil
}

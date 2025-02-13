package github

import (
	"fmt"
	"github_trending/internal/models"
)

type Service struct {
	client *Client
	parser *Parser
}

func NewService() *Service {
	return &Service{
		client: NewClient(),
		parser: NewParser(),
	}
}

func (s *Service) GetTrending(options models.TrendingOptions) ([]*models.Repository, error) {
	body, err := s.client.GetTrendingPage(options)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending page: %w", err)
	}
	defer body.Close()

	repositories, err := s.parser.ParseTrendingPage(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trending page: %w", err)
	}

	return repositories, nil
}

package github

import (
	"fmt"
	"github_trending/internal/models"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	client  *http.Client
	baseURL string
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: "https://github.com",
	}
}

func (c *Client) GetTrendingPage(options models.TrendingOptions) (io.ReadCloser, error) {
	url := c.buildURL(options)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending page: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}

func (c *Client) buildURL(options models.TrendingOptions) string {
	url := c.baseURL + "/trending"
	if options.Language != "" {
		url += "/" + options.Language
	}

	params := make([]string, 0)
	if options.Since != "" {
		params = append(params, "since="+options.Since)
	}
	if options.SpokenLanguage != "" {
		params = append(params, "spoken_language_code="+options.SpokenLanguage)
	}

	if len(params) > 0 {
		url += "?" + strings.Join(params, "&")
	}

	return url
}

package github

import (
	"fmt"
	"github_trending/internal/models"
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseTrendingPage(body io.Reader) ([]*models.Repository, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	var repositories []*models.Repository
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		repo := p.parseRepository(s)
		repositories = append(repositories, repo)
	})

	return repositories, nil
}

func (p *Parser) parseRepository(s *goquery.Selection) *models.Repository {
	repo := &models.Repository{}

	// Parse title and author
	if titleSel := s.Find("h2.h3 a"); titleSel.Length() > 0 {
		if repoPath, exists := titleSel.Attr("href"); exists {
			parts := strings.Split(strings.Trim(repoPath, "/"), "/")
			if len(parts) == 2 {
				repo.Author = parts[0]
				repo.Name = parts[1]
				repo.URL = "https://github.com" + repoPath
			}
		}
	}

	// Parse description
	repo.Description = strings.TrimSpace(s.Find("p").Text())

	// Parse language and color
	if langSel := s.Find("[itemprop='programmingLanguage']"); langSel.Length() > 0 {
		repo.Language = strings.TrimSpace(langSel.Text())
		if colorStyle, exists := langSel.Prev().Attr("style"); exists {
			if strings.Contains(colorStyle, "background-color") {
				parts := strings.Split(colorStyle, ":")
				if len(parts) == 2 {
					repo.LanguageColor = strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// Parse statistics
	if statsSel := s.Find("div.f6"); statsSel.Length() > 0 {
		repo.Stars = p.parseNumber(statsSel.Find("a[href$='/stargazers']").Text())
		repo.Forks = p.parseNumber(statsSel.Find("a[href$='/forks']").Text())
		repo.StarsGained = p.parseStarsGained(s)
	}

	return repo
}

func (p *Parser) parseNumber(text string) int {
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	num, err := strconv.Atoi(text)
	if err != nil {
		return 0 // Trả về 0 nếu không parse được số
	}
	return num
}

func (p *Parser) parseStarsGained(s *goquery.Selection) int {
	text := s.Find("span.d-inline-block.float-sm-right").Text()
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, ",", "")
	parts := strings.Split(text, " ")
	if len(parts) > 0 {
		num, _ := strconv.Atoi(parts[0])
		return num
	}
	return 0
}

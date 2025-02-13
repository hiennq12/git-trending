package models

type Repository struct {
	Author        string `json:"author"`
	Name          string `json:"name"`
	URL           string `json:"url"`
	Description   string `json:"description"`
	Language      string `json:"language"`
	LanguageColor string `json:"languageColor"`
	Stars         int    `json:"stars"`
	Forks         int    `json:"forks"`
	StarsGained   int    `json:"starsGained"`
}

type TrendingOptions struct {
	Language       string
	SpokenLanguage string
	Since          string
}

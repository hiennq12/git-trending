package telegram

import (
	"fmt"
	"github_trending/internal/models"
	"github_trending/pkg/utils"
	"strings"
)

func BuildMessage(repos []*models.Repository) string {
	message := strings.Builder{}
	message.WriteString(fmt.Sprintf(" %v \n", utils.GetToday()))

	for i, repo := range repos {
		message.WriteString(fmt.Sprintf("%v. Name: %v \n", i+1, repo.Name))
		message.WriteString(fmt.Sprintf("\t Language: %v \n", repo.Language))
		message.WriteString(fmt.Sprintf("\t Description: %v \n", repo.Description))
		message.WriteString(fmt.Sprintf("\t URL: %v \n\n", repo.URL))
	}

	return message.String()
}

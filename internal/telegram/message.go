package telegram

import (
	"fmt"
	"github_trending/internal/models"
	"github_trending/pkg/utils"
	"strings"
)

func BuildMessage(repo *models.Repository, rank int) string {
	message := strings.Builder{}
	message.WriteString(fmt.Sprintf(" %v \n", utils.GetToday()))
	message.WriteString(fmt.Sprintf("%v. Name: %v \n", rank+1, repo.Name))
	message.WriteString(fmt.Sprintf("\t Language: %v \n", repo.Language))
	message.WriteString(fmt.Sprintf("\t Original Description: %v \n", repo.Description))
	message.WriteString("\t Detailed Analysis:\n")
	message.WriteString(fmt.Sprintf("%v \n", repo.EnhancedDescription))
	message.WriteString(fmt.Sprintf("\t URL: %v \n\n", repo.URL))
	return message.String()
}

package main

import (
	"fmt"
	"github_trending/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// Repository represents a GitHub trending repository
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

// TrendingOptions contains parameters for fetching trending repositories
type TrendingOptions struct {
	Language       string // Programming language filter
	SpokenLanguage string // Spoken language filter
	Since          string // Time range: daily, weekly, monthly
}

// GetTrending fetches trending repositories from GitHub
func GetTrending(options TrendingOptions) ([]*Repository, error) {
	// Build URL with query parameters
	url := "https://github.com/trending"
	if options.Language != "" {
		url += "/" + options.Language
	}

	query := "?"
	if options.Since != "" {
		query += "since=" + options.Since
	}
	if options.SpokenLanguage != "" {
		if query != "?" {
			query += "&"
		}
		query += "spoken_language_code=" + options.SpokenLanguage
	}
	if query != "?" {
		url += query
	}

	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trending page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}

	// Parse HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	repositories := make([]*Repository, 0)

	// Find and parse each repository
	doc.Find("article.Box-row").Each(func(i int, s *goquery.Selection) {
		repo := &Repository{}

		// Get repository name and author
		titleSel := s.Find("h2.h3 a")
		if repoPath, exists := titleSel.Attr("href"); exists {
			parts := strings.Split(strings.Trim(repoPath, "/"), "/")
			if len(parts) == 2 {
				repo.Author = parts[0]
				repo.Name = parts[1]
				repo.URL = "https://github.com" + repoPath
			}
		}

		// Get description
		repo.Description = strings.TrimSpace(s.Find("p").Text())

		// Get programming language
		langSel := s.Find("[itemprop='programmingLanguage']")
		repo.Language = strings.TrimSpace(langSel.Text())
		if colorStyle, exists := langSel.Prev().Attr("style"); exists {
			if strings.Contains(colorStyle, "background-color") {
				parts := strings.Split(colorStyle, ":")
				if len(parts) == 2 {
					repo.LanguageColor = strings.TrimSpace(parts[1])
				}
			}
		}

		// Get statistics
		statsSel := s.Find("div.f6")

		// Parse stars
		starText := strings.TrimSpace(statsSel.Find("a[href$='/stargazers']").Text())
		starText = strings.ReplaceAll(starText, ",", "")
		if stars, err := strconv.Atoi(starText); err == nil {
			repo.Stars = stars
		}

		// Parse forks
		forkText := strings.TrimSpace(statsSel.Find("a[href$='/forks']").Text())
		forkText = strings.ReplaceAll(forkText, ",", "")
		if forks, err := strconv.Atoi(forkText); err == nil {
			repo.Forks = forks
		}

		// Parse stars gained
		starsGainedText := strings.TrimSpace(s.Find("span.d-inline-block.float-sm-right").Text())
		starsGainedText = strings.ReplaceAll(starsGainedText, ",", "")
		starsGainedText = strings.Split(starsGainedText, " ")[0]
		if starsGained, err := strconv.Atoi(starsGainedText); err == nil {
			repo.StarsGained = starsGained
		}

		repositories = append(repositories, repo)
	})

	return repositories, nil
}

// Example usage
func main() {
	// Get all trending repositories
	allTrending, err := GetTrending(TrendingOptions{})
	if err != nil {
		fmt.Printf("Error getting trending repos: %v\n", err)
		return
	}

	message := buildTextMessageSendTelegram(allTrending)
	//// Get JavaScript trending repositories for the week
	//jsTrending, err := GetTrending(TrendingOptions{
	//	Language: "javascript",
	//	Since:    "weekly",
	//})
	//if err != nil {
	//	fmt.Printf("Error getting JS trending repos: %v\n", err)
	//	return
	//}
	//fmt.Printf("Found %d JavaScript trending repositories\n", len(jsTrending))
	//
	//// Get Python trending repositories with English interface
	//pythonTrending, err := GetTrending(TrendingOptions{
	//	Language:       "python",
	//	SpokenLanguage: "en",
	//})
	//if err != nil {
	//	fmt.Printf("Error getting Python trending repos: %v\n", err)
	//	return
	//}
	//fmt.Printf("Found %d Python trending repositories\n", len(pythonTrending))

	// =========================================================
	// =========================================================
	// =======SEND MESSAGE======================================
	// =========================================================
	// =========================================================

	// Khởi tạo bot với token của bạn
	// token:
	bot, err := tgbotapi.NewBotAPI("TOKEN")
	if err != nil {
		log.Panic(err)
	}
	// chat id -4703770784
	chatID := int64(123123123)
	sendMessage(bot, chatID, message)
	log.Println("Đã gửi tin nhắn thành công!")
	//bot.Debug = true
	//log.Printf("Đã authorized vào tài khoản %s", bot.Self.UserName)
	//
	//// Cấu hình để nhận updates
	//u := tgbotapi.NewUpdate(0)
	//u.Timeout = 60
	//
	//// Tạo channel để nhận updates
	//updates, err := bot.GetUpdatesChan(u)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//// Xử lý tin nhắn đến
	//for update := range updates {
	//	if update.Message == nil {
	//		continue
	//	}
	//
	//	// Hiển thị tin nhắn nhận được
	//	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//
	//	// Tạo tin nhắn phản hồi
	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Bạn đã gửi: "+update.Message.Text)
	//	msg.ReplyToMessageID = update.Message.MessageID
	//
	//	// Gửi tin nhắn
	//	if _, err := bot.Send(msg); err != nil {
	//		log.Panic(err)
	//	}
	//}

	// =========================================================
	// =========================================================
	// =====CRON JOB ===========================================
	// =========================================================
	// =========================================================

}

func buildTextMessageSendTelegram(listRepo []*Repository) string {
	message := strings.Builder{}
	message.WriteString(fmt.Sprintf(" %v \n", utils.GetToday()))
	for i, repo := range listRepo {
		message.WriteString(fmt.Sprintf("%v. Name: %v \n", i+1, repo.Name))
		message.WriteString(fmt.Sprintf("\t Language: %v \n", repo.Language))
		message.WriteString(fmt.Sprintf("\t Description: %v \n", repo.Description))
		message.WriteString(fmt.Sprintf("\t URL: %v \n\n", repo.URL))
	}
	return message.String()
}

// =================================================================
// =================================================================
// ====SEND MESSAGE=================================================
// =================================================================
// =================================================================

// Hàm gửi tin nhắn đơn giản
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Println("Lỗi khi gửi tin nhắn:", err)
	}
}

// Hàm gửi ảnh
func sendPhoto(bot *tgbotapi.BotAPI, chatID int64, photoPath string, caption string) {
	photo := tgbotapi.NewPhotoUpload(chatID, photoPath)
	photo.Caption = caption
	if _, err := bot.Send(photo); err != nil {
		log.Println("Lỗi khi gửi ảnh:", err)
	}
}

// Hàm gửi document
func sendDocument(bot *tgbotapi.BotAPI, chatID int64, filePath string) {
	doc := tgbotapi.NewDocumentUpload(chatID, filePath)
	if _, err := bot.Send(doc); err != nil {
		log.Println("Lỗi khi gửi file:", err)
	}
}

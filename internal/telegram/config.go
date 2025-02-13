package telegram

const (
	DefaultBotToken = "TOKEN"
	DefaultChatID   = int64(-123456789)
)

type Config struct {
	BotToken string
	ChatID   int64
}

func NewConfig() *Config {
	return &Config{
		BotToken: DefaultBotToken,
		ChatID:   DefaultChatID,
	}
}

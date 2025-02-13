package telegram

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
)

const (
	DefaultBotToken = "TOKEN"
	DefaultChatID   = int64(-123456789)
)

type Config struct {
	TelegramConfig TelegramConfig `yaml:"telegram_config"`
}

type TelegramConfig struct {
	BotToken string `yaml:"bot_token"`
	ChatID   int64  `yaml:"chat_id"`
} // `yaml:"telegram_config"`

func NewConfig() *TelegramConfig {
	// Lấy đường dẫn của file hiện tại (config.go)
	_, filename, _, _ := runtime.Caller(0)
	// Lấy đường dẫn đến thư mục gốc của project
	projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	// Tạo đường dẫn đến file config.yaml
	configPath := filepath.Join(projectRoot, "config.yaml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := &Config{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	return &TelegramConfig{
		config.TelegramConfig.BotToken,
		config.TelegramConfig.ChatID,
	}
}

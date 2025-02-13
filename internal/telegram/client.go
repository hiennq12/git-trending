package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Client struct {
	bot    *tgbotapi.BotAPI
	config *Config
}

func NewClient(config *Config) (*Client, error) {
	bot, err := tgbotapi.NewBotAPI(config.BotToken)
	if err != nil {
		return nil, err
	}

	return &Client{
		bot:    bot,
		config: config,
	}, nil
}

func (c *Client) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(c.config.ChatID, text)
	_, err := c.bot.Send(msg)
	if err != nil {
		return err
	}
	return nil
}

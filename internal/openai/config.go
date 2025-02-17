package openai

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	APIKey string `yaml:"api_key"`
}

type OpenAIConfig struct {
	OpenAIConfig Config `yaml:"openai_config"`
}

func NewConfig() *Config {
	// Sử dụng cùng cách đọc config như telegram
	configPath := "/path/to/config.yaml" // Điều chỉnh path tương tự như trong telegram/config.go

	data, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	config := &OpenAIConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		panic(err)
	}
	return &Config{
		APIKey: config.OpenAIConfig.APIKey,
	}
}

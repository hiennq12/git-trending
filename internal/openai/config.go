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
	//// Lấy đường dẫn của file hiện tại (config.go)
	//_, filename, _, _ := runtime.Caller(0)
	//// Lấy đường dẫn đến thư mục gốc của project
	//projectRoot := filepath.Join(filepath.Dir(filename), "..", "..")
	//// Tạo đường dẫn đến file config.yaml
	//configPath := filepath.Join(projectRoot, "config.yaml")

	// config in server
	configPath := "/opt/github_trending/config/config.yaml"
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

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken    string
	DatabaseURL string
	AdminChat   string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("не удалось загрузить .env файл: %v", err)
	}

	cfg := &Config{
		BotToken:    os.Getenv("BOT_TOKEN"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		AdminChat:   os.Getenv("ADMINID"),
	}

	// Проверяем наличие обязательных переменных окружения
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("не задан BOT_TOKEN")
	}
	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("не задан DATABASE_URL")
	}
	if cfg.AdminChat == "" {
		return nil, fmt.Errorf("не задан ADMINID")
	}

	return cfg, nil
}

package bot

import (
	"fmt"
	"log"

	"github.com/SobolevTim/secret-santa-bot/internal/cache"
	"github.com/SobolevTim/secret-santa-bot/internal/database"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

type Bot struct {
	Client *telego.Bot
}

func NewBot(token string) (*Bot, error) {
	bot, err := telego.NewBot(token)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запуске бота: %v", err)
	}
	return &Bot{Client: bot}, nil
}

// Запуск бота
func (b *Bot) Start(service *database.Service, cache *cache.Cache) error {
	updates, err := b.Client.UpdatesViaLongPolling(nil)
	if err != nil {
		return fmt.Errorf("ошибка при получении обновлений от telegram: %v", err)
	}

	for update := range updates {
		if update.Message != nil {
			go b.handleMessage(update.Message, service, cache)
		}
	}
	return nil
}

func (b *Bot) SendMessage(chatID int64, msg string) {
	if _, err := b.Client.SendMessage(tu.Message(tu.ID(chatID), msg)); err != nil {
		log.Printf("ERROR: Ошибка отправки сообщения: %v", err)
	}
}

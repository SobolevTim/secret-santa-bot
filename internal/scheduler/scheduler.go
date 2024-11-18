package scheduler

import (
	"log"
	"time"

	"github.com/SobolevTim/secret-santa-bot/internal/bot"
	"github.com/SobolevTim/secret-santa-bot/internal/cache"
)

func StartDailyCacheClearer(c *cache.Cache, b *bot.Bot) {
	go func() {
		for {
			// Ждем до следующего 00:00
			now := time.Now()
			nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
			time.Sleep(time.Until(nextMidnight))

			// Уведомляем всех пользователей перед очисткой кеша
			notifyAndClearCache(c, b)
		}
	}()
}

func notifyAndClearCache(c *cache.Cache, b *bot.Bot) {
	// Получаем все пользовательские идентификаторы
	userIDs := c.Keys()

	for _, userID := range userIDs {
		// Отправляем сообщение каждому пользователю
		b.SendMessage(userID, "Ваше состояние было сброшено. Если вы не завершили действие, начните заново с команды /start.")
	}

	// Очищаем кеш
	c.Clear()
	log.Println("Cache cleared at midnight.")
}

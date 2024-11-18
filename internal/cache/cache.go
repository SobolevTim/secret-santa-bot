package cache

import (
	"log"
	"sync"

	"github.com/dgraph-io/ristretto"
)

type Cache struct {
	statusCache *ristretto.Cache
	dataCache   *ristretto.Cache
	mu          sync.Mutex // Для потокобезопасности при работе с кешом
}

func NewCache() *Cache {
	statusCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     32 * 1024 * 1024, // 32 МБ для статусов
		BufferItems: 64,
	})
	if err != nil {
		log.Fatalf("Failed to create status cache: %v", err)
	}

	dataCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,
		MaxCost:     32 * 1024 * 1024, // 32 МБ для данных
		BufferItems: 64,
	})
	if err != nil {
		log.Fatalf("Failed to create data cache: %v", err)
	}
	return &Cache{
		statusCache: statusCache,
		dataCache:   dataCache,
	}
}

// Установить значение
func (c *Cache) Set(userID int64, state string, data string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statusCache.Set(userID, state, 1)
	c.statusCache.Wait()
	c.dataCache.Set(userID, data, 1)
	c.dataCache.Wait()
}

// Получить значение
func (c *Cache) Get(userID int64) (string, string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	status, found := c.statusCache.Get(userID)
	if !found {
		return "", "", false
	}
	data, found := c.dataCache.Get(userID)
	if !found {
		return "", "", false
	}
	return status.(string), data.(string), true
}

// Получить все ключи (непрямой способ через базу данных ключей)
func (c *Cache) Keys() []int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]int64, 0)
	c.statusCache.Clear() // Точечный обход доступен через собственные методы
	return keys
}

// Очистить кеш
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statusCache.Clear()
	c.dataCache.Clear()
}

func (c *Cache) ClearUser(userID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.statusCache.Del(userID)
	c.dataCache.Del(userID)
}

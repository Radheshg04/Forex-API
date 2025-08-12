package cache

import (
	"greedygame/models"
	"sync"
	"time"
)

type ExchangeRateCache struct {
	currentRates    models.ExchangeRateResponse
	currentUpdated  time.Time
	historicalRates map[string]models.ExchangeRateResponse
	mu              sync.RWMutex
}

var cache *ExchangeRateCache
var once sync.Once

func GetCache() *ExchangeRateCache {
	once.Do(func() {
		cache = &ExchangeRateCache{
			historicalRates: make(map[string]models.ExchangeRateResponse),
		}
	})
	return cache

}

func (c *ExchangeRateCache) ReadCurrentCache() (rates models.ExchangeRateResponse, exists bool) {
	c.mu.RLock()
	rates, exists = c.currentRates, !c.currentUpdated.IsZero()
	c.mu.RUnlock()
	return rates, exists
}

func (c *ExchangeRateCache) WriteCurrentCache(value models.ExchangeRateResponse) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.currentRates = value
	c.currentUpdated = time.Now()
	return nil
}

func (c *ExchangeRateCache) ReadHistoricalCache(date string) (rates models.ExchangeRateResponse, exists bool) {
	c.mu.RLock()
	rates, exists = c.historicalRates[date]
	c.mu.RUnlock()
	return rates, exists
}

func (c *ExchangeRateCache) WriteHistoricalCache(dateKey string, rates models.ExchangeRateResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.historicalRates[dateKey] = rates
}

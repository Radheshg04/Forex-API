package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"greedygame/models"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var rdb *redis.Client

func InitCache() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	fmt.Println("Connected to Redis!")
}

func WriteHistoricalData(dateKey string, rates models.ExchangeRateResponse) error {
	jsonData, err := json.Marshal(rates)
	if err != nil {
		return fmt.Errorf("error marshalling historical data: %w", err)
	}
	err = rdb.Set(ctx, dateKey, jsonData, 90*24*time.Hour).Err()
	if err != nil {
		log.Printf("Error setting key: %v", err)
	}
	return nil
}

func ReadHistoricalData(dateKey string) (models.ExchangeRateResponse, error) {
	var rates models.ExchangeRateResponse

	val, err := rdb.Get(ctx, dateKey).Result()
	if err == redis.Nil {
		return rates, fmt.Errorf("no data found for %s", dateKey)
	} else if err != nil {
		return rates, fmt.Errorf("error getting key: %w", err)
	}

	if err := json.Unmarshal([]byte(val), &rates); err != nil {
		return rates, fmt.Errorf("error unmarshalling data: %w", err)
	}

	return rates, nil
}

func WriteCurrentData(value models.ExchangeRateResponse) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshalling current data: %w", err)
	}
	err = rdb.Set(ctx, "Current", jsonData, time.Hour).Err()
	if err != nil {
		log.Printf("Error setting key: %v", err)
	}
	return nil
}

func ReadCurrentData() (models.ExchangeRateResponse, error) {
	var rates models.ExchangeRateResponse
	val, err := rdb.Get(ctx, "Current").Result()
	if err != nil {
		return models.ExchangeRateResponse{}, err
	}

	if err := json.Unmarshal([]byte(val), &rates); err != nil {
		return rates, fmt.Errorf("error unmarshalling data: %w", err)
	}
	return rates, nil
}

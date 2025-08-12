package models

import (
	"sync"
	"time"
)

type ExchangeRateResponse struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

type ExchangeHostResponse struct {
	Success         bool               `json:"success"`
	BaseCode        string             `json:"source"`
	ConversionRates map[string]float64 `json:"quotes"`
}

type ExchangeRateCache struct {
	CurrentRates    ExchangeRateResponse
	CurrentUpdated  time.Time
	HistoricalRates map[string]ExchangeRateResponse
	Mu              sync.RWMutex
}

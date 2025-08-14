package services

import (
	"encoding/json"
	"fmt"
	"greedygame/cache"
	"greedygame/metrics"
	"greedygame/models"
	"log"
	"net/http"
	"os"
)

func UpdateCurrentForexInCache() error {
	var data models.ExchangeRateResponse
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", os.Getenv("EXCHANGE_RATE_API_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching exchange rates: ", err)
		return err
	}
	metrics.ExchangeRateApiCalls.Inc()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return err
	}
	if data.Result != "success" {
		// TODO: Call fallback
		data, err = ForexFallback()
		if err != nil {
			log.Println("Couldnt Update cache from fallback: ", err)
			return err
		}
	}
	filtered := make(map[string]float64)
	for k, v := range data.ConversionRates {
		if supportedCurrencies[k] {
			filtered[k] = v
		}
	}
	data.ConversionRates = filtered
	cache.GetCache().WriteCurrentCache(data)
	if err != nil {
		log.Println("Error updating cache:", err)
		return err
	}
	return nil
}

func UpdateHistoricalCache(year, month, day int) error {
	var data models.ExchangeRateResponse
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/history/USD/%d/%d/%d", os.Getenv("EXCHANGE_RATE_API_KEY"), year, month, day)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching latest exchange rates, ", err.Error())
		return err
	}
	metrics.ExchangeRateApiCalls.Inc()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return err
	}
	if data.Result != "success" {
		data, err = HistoricalForexFallback(year, month, day)
		metrics.FallbackExchangeRateApiCalls.Inc()
		if err != nil {
			log.Println("Error in historical forex fallback:", err)
			return err
		}
	}
	dateKey := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	cache.GetCache().WriteHistoricalCache(dateKey, data)
	return nil
}

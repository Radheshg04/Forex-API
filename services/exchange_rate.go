package services

import (
	"encoding/json"
	"fmt"
	"greedygame/cache"
	"greedygame/models"
	"log"
	"net/http"
	"os"
)

var supportedCurrencies = map[string]bool{
	"USD": true,
	"INR": true,
	"EUR": true,
	"JPY": true,
	"GBP": true,
}

// This function standardizes the fallback API in main API format for further processing
func standardizeAPI(data models.ExchangeHostResponse) models.ExchangeRateResponse {
	var response models.ExchangeRateResponse
	newMap := make(map[string]float64)
	for i, j := range data.ConversionRates {
		currency := i[3:]
		if supportedCurrencies[currency] {
			newMap[currency] = j
		}
	}
	response.ConversionRates = newMap
	if data.Success {
		response.Result = "success"
	} else {
		response.Result = "error"
	}
	response.BaseCode = data.BaseCode
	return response
}
func ForexFallback() (models.ExchangeRateResponse, error) {
	fmt.Println("Inside fallback function")
	var data models.ExchangeHostResponse
	var response models.ExchangeRateResponse
	url := fmt.Sprintf("https://api.exchangerate.host/live?access_key=%s", os.Getenv("EXCHANGE_HOST_API_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching latest exchange rates, ", err.Error())
		return response, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return response, err
	}
	response = standardizeAPI(data)
	return response, nil
}

func HistoricalForexFallback(year int, month int, day int) (models.ExchangeRateResponse, error) {
	var data models.ExchangeHostResponse
	var response models.ExchangeRateResponse
	url := fmt.Sprintf("https://api.exchangerate.host/historical?access_key=%s&date=%d-%d-%d", os.Getenv("EXCHANGE_HOST_API_KEY"), year, month, day)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching latest exchange rates, ", err.Error())
		return response, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return response, err
	}
	response = standardizeAPI(data)
	return response, nil
}

// Function for fetching exchange rates
func ListExchangeRates() models.ExchangeRateResponse {
	rates, exists := cache.GetCache().ReadCurrentCache()
	if !exists {
		log.Println("Cache miss!")
		return models.ExchangeRateResponse{}
	}
	return rates
}

// Function to get historical exchange rates
func ListHistoricalExchangeRates(year int, month int, day int) (models.ExchangeRateResponse, error) {
	rates, exists := cache.GetCache().ReadHistoricalCache(fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if !exists {
		return rates, fmt.Errorf("records older than 90 days are not persisted")
	}
	return rates, nil
}

func GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	rates, exists := cache.GetCache().ReadCurrentCache()
	if !exists {
		log.Printf("Cache not initialized")
	}
	if fromCurrency == toCurrency {
		return 1.0, nil
	}

	fromRate, fromExists := rates.ConversionRates[fromCurrency]
	toRate, toExists := rates.ConversionRates[toCurrency]

	if !fromExists {
		return 0, fmt.Errorf("rate not available for currency %s", fromCurrency)
	}
	if !toExists {
		return 0, fmt.Errorf("rate not available for currency %s", toCurrency)
	}

	if fromCurrency == "USD" {
		return toRate, nil
	}
	if toCurrency == "USD" {
		return 1.0 / fromRate, nil
	}
	return toRate / fromRate, nil
}

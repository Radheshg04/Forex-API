package services

import (
	"encoding/json"
	"fmt"
	"greedygame/models"
	"log"
	"net/http"
	"os"
)

// This function standardizes the fallback API in main API format for further processing
func standardizeAPI(data models.ExchangeHostResponse) models.ExchangeRateResponse {
	var response models.ExchangeRateResponse
	newMap := make(map[string]float64)
	for k, v := range data.ConversionRates {
		newMap[k[3:]] = v
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

// Function for fetching exchange rates
func GetExchangeRates() (models.ExchangeRateResponse, error) {
	var data models.ExchangeRateResponse
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/latest/USD", os.Getenv("EXCHANGE_RATE_API_KEY"))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching exchange rates: ", err)
		return data, nil
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return data, err
	}
	if data.Result != "success" {
		return GetExchangeRatesFallback()
	}
	return data, nil
}

// Fallback function for fetching exchange rates
// Used if first API is down
func GetExchangeRatesFallback() (models.ExchangeRateResponse, error) {
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

// Function to get historical exchange rates
func GetHistoricalExchangeRates(year int, month int, day int) (models.ExchangeRateResponse, error) {
	var data models.ExchangeRateResponse
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/%s/history/USD/%d/%d/%d", os.Getenv("EXCHANGE_RATE_API_KEY"), year, month, day)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching latest exchange rates, ", err.Error())
		return data, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding JSON:", err)
		return data, err
	}
	if data.Result != "success" {
		return GetHistoricalExchangeRatesFallback(year, month, day)
	}
	return data, nil
}

// Fallback function for fetching historical exchange rates
// Used if first API is down
func GetHistoricalExchangeRatesFallback(year int, month int, day int) (models.ExchangeRateResponse, error) {
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

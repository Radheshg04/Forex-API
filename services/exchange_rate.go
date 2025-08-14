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
	"sync"
	"time"
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
	metrics.FallbackExchangeRateApiCalls.Inc()
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
	metrics.FallbackExchangeRateApiCalls.Inc()
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
		metrics.CacheMisses.Inc()
		log.Println("Cache miss!")
		return models.ExchangeRateResponse{}
	}
	metrics.CacheHits.Inc()
	return rates
}

// Function to get historical exchange rates
func ListHistoricalExchangeRates(year int, month int, day int) (models.ExchangeRateResponse, error) {
	rates, exists := cache.GetCache().ReadHistoricalCache(fmt.Sprintf("%d-%02d-%02d", year, month, day))
	if !exists {
		return rates, fmt.Errorf("records older than 90 days are not persisted")
	}
	metrics.CacheHits.Inc()
	return rates, nil
}

func GetCurrentExchangeRate(fromCurrency, toCurrency string) (float64, error) {
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
	metrics.CacheHits.Inc()
	if !toExists {
		return 0, fmt.Errorf("rate not available for currency %s", toCurrency)
	}
	metrics.CacheHits.Inc()

	if fromCurrency == "USD" {
		return toRate, nil
	}
	if toCurrency == "USD" {
		return 1.0 / fromRate, nil
	}
	return toRate / fromRate, nil
}

func GetHistoricalExchangeRate(fromCurrency, toCurrency string, date string) (float64, error) {
	requestedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if !supportedCurrencies[fromCurrency] || !supportedCurrencies[toCurrency] {
		log.Printf("Invalid Request")
		return 0, fmt.Errorf("invalid or unsupported currencies")
	}

	threshold := time.Now().AddDate(0, 0, -90)
	if requestedDate.Before(threshold) {
		log.Printf("Data requested older than 90 days. Requested date: %v", requestedDate)
		return 0, fmt.Errorf("historical data older than 90 days cannot be retrieved")
	}
	threshold = time.Now()
	if threshold.Before(requestedDate) {
		log.Printf("Future data. Requested date: %v", requestedDate)
		return 0, fmt.Errorf("requested date is in the future")
	}
	if fromCurrency == toCurrency {
		return 1.0, nil
	}
	rates, exists := cache.GetCache().ReadHistoricalCache(date)
	if !exists {
		log.Printf("Cache miss")
		metrics.CacheMisses.Inc()
		err = ensureSingleCacheUpdate(date)
		if err != nil {
			return 0, err
		}
		rates, _ = cache.GetCache().ReadHistoricalCache(date)
	}
	metrics.CacheHits.Inc()

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

func GetHistoricalExchangeRatesOverTimeRange(fromCurrency, toCurrency string, fromDate, toDate string) ([]float64, error) {
	requestedFromDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	requestedToDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD")
	}
	if !supportedCurrencies[fromCurrency] || !supportedCurrencies[toCurrency] {
		log.Printf("Invalid Request: unsupported currencies")
		return nil, fmt.Errorf("invalid or unsupported currencies")
	}
	if requestedToDate.Before(requestedFromDate) {
		return nil, fmt.Errorf("invalid date range")
	}

	if time.Now().Before(requestedFromDate) || time.Now().Before(requestedToDate) {
		log.Printf("Future data requested.")
		return nil, fmt.Errorf("requested date is in the future")
	}
	ninetyDaysAgo := time.Now().AddDate(0, 0, -90)
	if requestedFromDate.Before(ninetyDaysAgo) || requestedToDate.Before(ninetyDaysAgo) {
		log.Printf("Data requested older than 90 days.")
		return nil, fmt.Errorf("historical data older than 90 days cannot be retrieved")
	}

	var exchangeRates []float64

	for d := requestedFromDate; !d.After(requestedToDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")

		rates, exists := cache.GetCache().ReadHistoricalCache(dateStr)
		if !exists {
			metrics.CacheMisses.Inc()
			err = ensureSingleCacheUpdate(dateStr)
			if err != nil {
				return nil, err
			}
			rates, exists = cache.GetCache().ReadHistoricalCache(dateStr)
			if !exists {
				return nil, fmt.Errorf("record not found")
			}
		}
		metrics.CacheHits.Inc()
		fromRate := rates.ConversionRates[fromCurrency]
		toRate := rates.ConversionRates[toCurrency]

		var calculatedRate float64
		if fromCurrency == "USD" {
			calculatedRate = toRate
		} else if toCurrency == "USD" {
			calculatedRate = 1.0 / fromRate
		} else {
			calculatedRate = toRate / fromRate
		}

		exchangeRates = append(exchangeRates, calculatedRate)
	}

	return exchangeRates, nil
}

func SingleFetcher() *models.SingleFetch {
	return &models.SingleFetch{
		Requests: make(map[string]*sync.WaitGroup),
	}
}

var sf = SingleFetcher()

func ensureSingleCacheUpdate(date string) error {
	sf.Mu.Lock()
	wg, exists := sf.Requests[date]

	if !exists {
		wg = new(sync.WaitGroup)
		wg.Add(1)
		sf.Requests[date] = wg
		sf.Mu.Unlock()

		requestedDate, _ := time.Parse("2006-01-02", date)
		err := UpdateHistoricalCache(requestedDate.Year(), int(requestedDate.Month()), requestedDate.Day())
		wg.Done()
		sf.Mu.Lock()
		delete(sf.Requests, date)
		sf.Mu.Unlock()

		if err != nil {
			return err
		}

	} else {
		sf.Mu.Unlock()
		wg.Wait()
	}
	return nil
}

package handlers

import (
	"greedygame/metrics"
	"greedygame/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ListExchangeRatesHandler(c *gin.Context) {
	metrics.ListExchangeRateRequests.Inc()
	start := time.Now()
	defer metrics.ListExchangeRateRequestDuration.Observe(time.Since(start).Seconds())
	data := services.ListExchangeRates()
	if data.Result == "" {
		metrics.ListExchangeRateErrors.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch Exchange Rates"})
		return
	}
	c.JSON(http.StatusOK, data.ConversionRates)
}

func GetForexOverTime(c *gin.Context) {
	metrics.GetForexOverTimeRequests.Inc()
	start := time.Now()
	defer metrics.GetForexOverTimeRequestDuration.Observe(time.Since(start).Seconds())

	currency1 := c.Query("currency1")
	currency2 := c.Query("currency2")
	fromDate := c.Query("from")
	tillDate := c.Query("to")
	if currency1 == "" || currency2 == "" || fromDate == "" || tillDate == "" {
		metrics.GetForexOverTimeErrors.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing Parameters"})
		return
	}
	data, err := services.GetHistoricalExchangeRatesOverTimeRange(currency1, currency2, fromDate, tillDate)
	if err != nil {
		metrics.GetForexOverTimeErrors.Inc()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversion rates": data})

}

func GetExchangeRateHandler(c *gin.Context) {
	metrics.GetExchangeRateRequests.Inc()
	start := time.Now()
	defer metrics.GetExchangeRateRequestDuration.Observe(time.Since(start).Seconds())

	from := c.Query("from")
	to := c.Query("to")
	date := c.Query("date")

	if from == "" || to == "" {
		metrics.GetExchangeRateErrors.Inc()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	if date != "" {
		exchangeRate, err := services.GetHistoricalExchangeRate(from, to, date)
		if err != nil {
			metrics.GetExchangeRateErrors.Inc()
			switch err.Error() {
			case "historical data older than 90 days cannot be retrieved", "requested date is in the future":
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, gin.H{"amount": exchangeRate})
	} else {
		exchangeRate, err := services.GetCurrentExchangeRate(from, to)
		if err != nil {
			metrics.GetExchangeRateErrors.Inc()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"amount": exchangeRate})
	}
}

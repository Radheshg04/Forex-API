package handlers

import (
	"greedygame/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ListExchangeRatesHandler(c *gin.Context) {
	data := services.ListExchangeRates()
	if data.Result == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch Exchange Rates"})
		return
	}
	c.JSON(http.StatusOK, data.ConversionRates)
}

func GetExchangeRatesOverTimePeriod(c *gin.Context) {
	currency1 := c.Query("currency1")
	currency2 := c.Query("currency2")
	fromDate := c.Query("from")
	tillDate := c.Query("to")
	data, err := services.GetHistoricalExchangeRatesOverTimeRange(currency1, currency2, fromDate, tillDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"conversion rates": data})

}

func GetExchangeRateHandler(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	date := c.Query("date")

	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	if date != "" {
		exchangeRate, err := services.GetHistoricalExchangeRate(from, to, date)
		if err != nil {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"amount": exchangeRate})
	}
}

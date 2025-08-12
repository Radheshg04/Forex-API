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
	c.JSON(http.StatusOK, data)
}

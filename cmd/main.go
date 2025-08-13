package main

import (
	"greedygame/handlers"
	"greedygame/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("No .env file found")
		return
	}
	services.UpdateCurrentForexInCache()
	go func() {
		services.PollFunction()
	}()
}

func main() {
	Init()

	r := gin.Default()
	r.GET("/rates/latest", handlers.ListExchangeRatesHandler)
	r.GET("/rates/historical", handlers.GetExchangeRatesOverTimePeriod)
	r.GET("/convert", handlers.GetExchangeRateHandler)
	r.Run(":8080")
}

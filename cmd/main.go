package main

import (
	"greedygame/handlers"
	"greedygame/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	r.GET("/rates/historical", handlers.GetForexOverTime)
	r.GET("/convert", handlers.GetExchangeRateHandler)

	metricsRouter := gin.New()
	metricsRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))
	go func() {
		metricsRouter.Run("0.0.0.0:9090")
	}()

	r.Run("0.0.0.0:8080")
}

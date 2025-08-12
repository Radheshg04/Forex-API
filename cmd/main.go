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
	services.PollFunction()

}

func main() {
	Init()
	// result, err := services.GetHistoricalExchangeRates(2020, 10, 2)
	// result, err := services.GetExchangeRates()
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }
	// fmt.Println(result)
	r := gin.Default()

	r.GET("/rates/latest", handlers.ListExchangeRatesHandler)
	// r.GET("/rates/historical/:year/:month/:day", handlers.GetHistoricalRatesHandler)
	// r.GET("/convert", handlers.ConvertCurrencyHandler)

	r.Run(":8080")
}

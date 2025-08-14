package services

import (
	"fmt"
	"greedygame/metrics"
	"log"
	"time"
)

// TODO: Add the update cache function
func PollFunction() {
	metrics.TotalPolls.Inc()
	metrics.LastPollTimestamp.Set(float64(time.Now().Unix()))
	pollingInterval := time.Hour
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			err := UpdateCurrentForexInCache()
			if err != nil {
				metrics.PollFails.Inc()
				log.Printf("Error fetching current forex rates: %s\n", err.Error())
			}
			metrics.TotalPolls.Inc()
			metrics.LastPollTimestamp.Set(float64(time.Now().Unix()))
			fmt.Printf("Exchange Rates updated at time %s", t.Format("15:04:05 IST"))
		}
	}
}

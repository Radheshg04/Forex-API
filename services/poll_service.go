package services

import (
	"fmt"
	"log"
	"time"
)

// TODO: Add the update cache function
func PollFunction() {
	pollingInterval := time.Hour
	ticker := time.NewTicker(pollingInterval)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			err := UpdateCurrentForexInCache()
			if err != nil {
				log.Printf("Error fetching current forex rates: %s\n", err.Error())
			}
			fmt.Printf("Exchange Rates updated at time %s", t.Format("15:04:05 IST"))
		}
	}
}

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Handler metrics
	GetExchangeRateRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_exchange_rate_requests_total",
		Help: "The total number of requests to the GetExchangeRateHandler.",
	})
	GetExchangeRateRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "get_exchange_rate_request_duration_seconds",
		Help:    "Histogram of the duration of requests to the GetExchangeRateHandler.",
		Buckets: prometheus.DefBuckets,
	})
	GetExchangeRateErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_exchange_rate_errors_total",
		Help: "The total number of errors returned by the GetExchangeRateHandler.",
	})

	ListExchangeRateRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "list_exchange_rate_requests_total",
		Help: "The total number of requests to the ListExchangeRateHandler.",
	})
	ListExchangeRateRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "list_exchange_rate_request_duration_seconds",
		Help:    "Histogram of the duration of requests to the ListExchangeRateHandler.",
		Buckets: prometheus.DefBuckets,
	})
	ListExchangeRateErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "list_exchange_rate_errors_total",
		Help: "The total number of errors returned by the ListExchangeRateHandler.",
	})

	GetForexOverTimeRequests = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_forex_over_time_requests_total",
		Help: "The total number of requests to the GetForexOverTimeHandler.",
	})
	GetForexOverTimeRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "get_forex_over_time_request_duration_seconds",
		Help:    "Histogram of the duration of requests to the GetForexOverTimeHandler.",
		Buckets: prometheus.DefBuckets,
	})
	GetForexOverTimeErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "get_forex_over_time_errors_total",
		Help: "The total number of errors returned by the GetForexOverTimeHandler.",
	})

	// Cache metrics
	CacheMisses = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_misses_total",
		Help: "The total number of cache hits.",
	})
	CacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "cache_hits_total",
		Help: "The total number of cache misses.",
	})

	// API call metrics
	ExchangeRateApiCalls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "exchange_rate_api_calls",
		Help: "Total number of times exchange rate api was called",
	})
	FallbackExchangeRateApiCalls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "fallback_exchange_rate_api_calls",
		Help: "Total number of times fallback exchange rate api was called",
	})

	// Polling Service Metrics
	TotalPolls = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_poll_count_total",
		Help: "Total number of times the polling service has polled.",
	})
	LastPollTimestamp = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "api_last_poll_timestamp_seconds",
		Help: "Unix timestamp of the last successful poll of the external API.",
	})
	PollFails = promauto.NewCounter(prometheus.CounterOpts{
		Name: "api_poll_fail_count_total",
		Help: "Total number of times the polling service failed to update.",
	})
)

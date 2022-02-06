package hw3

import "github.com/prometheus/client_golang/prometheus"

func RegisterMetrics() map[string]prometheus.Collector {
	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
		},
		[]string{"code", "endpoint"})
	prometheus.MustRegister(counter)

	histogram := prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "request_processing_time_histogram_ms",
		})
	prometheus.MustRegister(histogram)

	return map[string]prometheus.Collector{
		"counter":   counter,
		"histogram": histogram,
	}
}

package hw3

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"time"
)

type Monitor struct {
	requestsRate *prometheus.CounterVec
	latency      *prometheus.HistogramVec
	errorRate    *prometheus.CounterVec
}

var defaultBuckets = []float64{0.1, 0.3, 1.5, 10.5}

func NewMonitor() (*Monitor, error) {
	monitor := &Monitor{}

	monitor.requestsRate = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_number",
			Help: "HTTP requests rate.",
		},
		[]string{"method", "addr"},
	)

	monitor.latency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_latency",
			Help:    "Duration in seconds of HTTP requests.",
			Buckets: defaultBuckets,
		},
		[]string{"type", "status", "method", "addr", "errorMessage"},
	)

	monitor.errorRate = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_rate",
			Help: "Rate of responses with HTTP 500 status.",
		},
		[]string{"method", "addr"},
	)

	return monitor, nil
}

func (m *Monitor) Prometheus() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		started := time.Now()

		err = c.Next()

		duration := time.Since(started)

		statusCode := c.Response().StatusCode()

		m.requestsRate.WithLabelValues(c.Method(), c.Route().Path).Inc()
		m.latency.WithLabelValues(
			c.Protocol(),
			fmt.Sprint(statusCode),
			c.Method(),
			c.Route().Path,
			"",
			fmt.Sprintf("%v", duration.Seconds()),
		)

		return err
	}
}

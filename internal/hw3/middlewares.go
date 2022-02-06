package hw3

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *Service) WithMetricsMiddleware(counter *prometheus.CounterVec) *Service {
	s.Use(func(c *fiber.Ctx) error {
		counter.WithLabelValues(strconv.Itoa(c.Response().StatusCode()), c.OriginalURL()).Add(1)

		log.Printf("%d: %s", c.Response().StatusCode(), c.OriginalURL())

		return c.Next()
	})

	return s
}

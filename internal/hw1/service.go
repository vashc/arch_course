package hw1

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func NewService() *Service {
	return &Service{
		App: fiber.New(),
	}
}

func (s *Service) InstantiateRoutes() {
	s.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString(`{"status": "OK"}`)
	})

	s.Get("/student/:student", func(c *fiber.Ctx) error {
		return c.SendString(
			fmt.Sprintf("Hello, %s", c.Params("student", "unknown")),
		)
	})
}

func (s *Service) Start(port string) error {
	return s.Listen(port)
}

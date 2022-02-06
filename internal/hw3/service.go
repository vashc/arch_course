package hw3

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewService(config *Config, storage *Storage, monitor *Monitor) *Service {
	appConfig := fiber.Config{
		// Override default error handler
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Status code defaults to 500
			code := fiber.StatusInternalServerError

			// Retrieve the custom status code if it's an fiber.*Error
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			response, _ := json.Marshal(Response{
				Code:    code,
				Message: fmt.Sprintf("internal server error: %s", err.Error()),
			})

			monitor.errorRate.WithLabelValues(c.Method(), c.Route().Path).Inc()

			return c.Status(code).SendString(string(response))
		},
	}

	app := fiber.New(appConfig)
	app.Use(monitor.Prometheus())

	return &Service{
		config:  config,
		storage: storage,
		App:     app,
	}
}

func (s *Service) InstantiateRoutes() {
	s.Post("/user", s.createUserHandler())

	s.Get("/user/:userID", s.getUserHandler())

	s.Delete("/user/:userID", s.deleteUserHandler())

	s.Put("/user/:userID", s.updateUserHandler())

	s.Get("/health", s.healthHandler())

	s.Get("/metrics", s.metricsHandler())
}

func (s *Service) Start(port string) error {
	return s.Listen(port)
}

func (s *Service) createUserHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		log.Print("POST /user")

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		defer func() {
			code, response := s.getResponse(err, "user created")
			_ = c.Status(code).SendString(response)
		}()

		user := new(User)

		if err := c.BodyParser(user); err != nil {
			return err
		}

		if err := s.storage.CreateUser(user); err != nil {
			return err
		}

		return
	}
}

func (s *Service) getUserHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		log.Print("GET /user")

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		defer func() {
			code, response := s.getResponse(err, "user created")
			if code != fiber.StatusOK {
				log.Printf("getUser error: %s", err.Error())
				_ = c.Status(code).SendString(response)
			}
		}()

		userID, err := getUserIDFromCtx(c)
		if err != nil {
			return err
		}

		user, err := s.storage.GetUser(userID)
		if err != nil {
			return err
		}

		response, err := json.Marshal(user)
		if err != nil {
			return err
		}

		return c.SendString(string(response))
	}
}

func (s *Service) deleteUserHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		log.Print("DELETE /user")

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		defer func() {
			code, response := s.getResponse(err, "user deleted")
			_ = c.Status(code).SendString(response)
		}()

		userID, err := getUserIDFromCtx(c)
		if err != nil {
			return err
		}

		if err = s.storage.DeleteUser(userID); err != nil {
			return err
		}

		return
	}
}

func (s *Service) updateUserHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) (err error) {
		log.Print("PUT /user")

		c.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		defer func() {
			code, response := s.getResponse(err, "user updated")
			_ = c.Status(code).SendString(response)
		}()

		userID, err := getUserIDFromCtx(c)
		if err != nil {
			return err
		}

		user := &User{ID: userID}

		if err := c.BodyParser(user); err != nil {
			return err
		}

		if err := s.storage.UpdateUser(user); err != nil {
			return err
		}

		return
	}
}

func (s *Service) healthHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		log.Print("GET /health")

		response, err := json.Marshal(Response{
			Code:    200,
			Message: "OK",
		})
		if err != nil {
			return err
		}

		return c.SendString(string(response))
	}
}

func (s *Service) metricsHandler() func(*fiber.Ctx) error {
	log.Print("GET /metrics")
	return adaptor.HTTPHandler(promhttp.Handler())
}

func (s *Service) getResponse(err error, message string) (int, string) {
	code := fiber.StatusOK
	response, _ := json.Marshal(Response{
		Code:    code,
		Message: message,
	})

	if err != nil {
		code = fiber.StatusInternalServerError
		response, _ = json.Marshal(Response{
			Code:    code,
			Message: "internal server error",
		})
	}

	return code, string(response)
}

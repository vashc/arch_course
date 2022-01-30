package hw2

import (
	"github.com/gocraft/dbr/v2"
	"github.com/gofiber/fiber/v2"
)

type Service struct {
	config  *Config
	storage *Storage

	*fiber.App
}

type Storage struct {
	sess *dbr.Session
}

type Config struct {
	DBURI    string `envconfig:"DATABASE_URI" required:"true"`
	Hostname string `envconfig:"HOSTNAME"`
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

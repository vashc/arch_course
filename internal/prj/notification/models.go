package notification

import (
	"arch_course/internal/prj"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
)

type Service struct {
	config  *Config
	storage *Storage
	client  *prj.RabbitClient

	*chi.Mux
}

type Worker struct {
	config  *Config
	storage *Storage
	client  *prj.RabbitClient
}

type Storage struct {
	sess *dbr.Session
}

type Config struct {
	DBURI string `envconfig:"DATABASE_URI" required:"true"`

	RabbitConfig *prj.RabbitConfig
}

type Response struct {
	Status string `json:"status"`
}

package exchanger

import (
	"arch_course/internal/prj"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
	"net/http"
)

type Service struct {
	config       *Config
	storage      *Storage
	rabbitClient *prj.RabbitClient

	*chi.Mux
}

type Worker struct {
	config       *Config
	storage      *Storage
	client       *http.Client
	rabbitClient *prj.RabbitClient
}

type Storage struct {
	sess *dbr.Session
}

type Config struct {
	DBURI string `envconfig:"DATABASE_URI" required:"true"`

	BalanceHost string `envconfig:"BALANCE_HOST" required:"true"`
	BalancePort int    `envconfig:"BALANCE_PORT" required:"true"`

	RabbitConfig *prj.RabbitConfig
}

type Response struct {
	OrderUUID string `json:"order_uuid"`
	Status    string `json:"status"`
}

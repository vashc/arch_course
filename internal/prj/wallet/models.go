package wallet

import (
	"arch_course/internal/prj"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type Service struct {
	config       *Config
	storage      *Storage
	client       *http.Client
	rabbitClient *prj.RabbitClient

	*chi.Mux
}

type Storage struct {
	sess *dbr.Session
}

type Config struct {
	DBURI    string `envconfig:"DATABASE_URI" required:"true"`
	Hostname string `envconfig:"HOSTNAME"`

	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`

	BalanceHost string `envconfig:"BALANCE_HOST" required:"true"`
	BalancePort int    `envconfig:"BALANCE_PORT" required:"true"`

	NotificationHost string `envconfig:"NOTIFICATION_HOST" required:"true"`
	NotificationPort int    `envconfig:"NOTIFICATION_PORT" required:"true"`

	ExchangerHost string `envconfig:"EXCHANGER_HOST" required:"true"`
	ExchangerPort int    `envconfig:"EXCHANGER_PORT" required:"true"`

	BcgatewayHost string `envconfig:"BCGATEWAY_HOST" required:"true"`
	BcgatewayPort int    `envconfig:"BCGATEWAY_PORT" required:"true"`

	AuthHost string `envconfig:"AUTH_HOST" required:"true"`
	AuthPort int    `envconfig:"AUTH_PORT" required:"true"`

	RabbitConfig *prj.RabbitConfig
}

type Worker struct {
	config       *Config
	storage      *Storage
	client       *http.Client
	rabbitClient *prj.RabbitClient
}

type Order struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	Type         string  `json:"type"`
	CryptoAmount float64 `json:"crypto_amount"`
	FiatAmount   float64 `json:"fiat_amount"`
	Status       string  `json:"status"`
}

type OrderProcessing struct {
	ID          int64 `json:"id"`
	OrderID     int64 `json:"order_id"`
	StepsNumber int   `json:"steps_number"`
	FailedSteps []int `json:"failed_steps"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	OrderID int64  `json:"order_id,omitempty"`
	Status  string `json:"status"`
}

type JWTClaims struct {
	UserID int64 `json:"user_id"`

	jwt.StandardClaims
}

type AuthToken map[string]interface{}

package auth

import (
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
	"github.com/golang-jwt/jwt"
)

type Service struct {
	config  *Config
	storage *Storage

	*chi.Mux
}

type Storage struct {
	sess *dbr.Session
}

type Config struct {
	DBURI    string `envconfig:"DATABASE_URI" required:"true"`
	Hostname string `envconfig:"HOSTNAME"`

	JWTSecret string `envconfig:"JWT_SECRET" required:"true"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Status string `json:"status"`
}

type JWTClaims struct {
	UserID int64 `json:"user_id"`

	jwt.StandardClaims
}

type AuthToken map[string]interface{}

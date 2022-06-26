package billing_service

import (
	"arch_course/internal/hw6"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
)

type Service struct {
	Config  *hw6.Config
	Storage *Storage

	*chi.Mux
}

type Worker struct {
	Config  *hw6.Config
	Storage *Storage
	Client  *hw6.Client
}

type Storage struct {
	Sess *dbr.Session
}

type Account struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Balance  float64 `json:"balance"`
}

type BalanceRequest struct {
	Amount float64 `json:"amount"`
}

type Response struct {
	Status string `json:"status"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
}

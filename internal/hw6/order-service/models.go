package order_service

import (
	"arch_course/internal/hw6"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
	"time"
)

type Service struct {
	Config  *hw6.Config
	Storage *Storage
	Client  *hw6.Client

	*chi.Mux
}

type Storage struct {
	Sess *dbr.Session
}

type Order struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Status string `json:"status"`
}

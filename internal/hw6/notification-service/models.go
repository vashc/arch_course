package notification_service

import (
	"arch_course/internal/hw6"
	"github.com/go-chi/chi/v5"
	"github.com/gocraft/dbr/v2"
	"time"
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

type Message struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
}

type Response struct {
	Status string `json:"status"`
}

type MessagesResponse struct {
	Messages []Message `json:"messages"`
	Count    int       `json:"count"`
}

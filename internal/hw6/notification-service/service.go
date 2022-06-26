package notification_service

import (
	"arch_course/internal/hw6"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewService(config *hw6.Config, storage *Storage) *Service {
	return &Service{
		Config:  config,
		Storage: storage,
		Mux:     chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	s.Route("/message", func(router chi.Router) {
		router.Post("/", s.sendMessageHandler())
		router.Get("/{email}", s.getMessagesHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) sendMessageHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		message := new(Message)

		err := hw6.BodyParser(w, r, message)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		if err = s.Storage.SendMessage(message); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(Response{Status: http.StatusText(http.StatusOK)})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) getMessagesHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get email
		email := chi.URLParam(r, hw6.Email)
		if email == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		messages, err := s.Storage.GetMessagesByEmail(email)
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(MessagesResponse{
			Messages: messages,
			Count:    len(messages),
		})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) healthHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := json.Marshal(Response{Status: http.StatusText(http.StatusOK)})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

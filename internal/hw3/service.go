package hw3

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewService(config *Config, storage *Storage) *Service {
	return &Service{
		config:  config,
		storage: storage,
		Mux:     chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	s.Route("/user", func(router chi.Router) {
		router.Use(NewPatternMiddleware("hw3"))
		router.Get("/{user_id}", s.getUserHandler())
		router.Post("/", s.createUserHandler())
		router.Put("/{user_id}", s.updateUserHandler())
		router.Delete("/{user_id}", s.deleteUserHandler())
	})

	s.Get("/health", s.healthHandler())

	s.Get("/metrics", s.metricsHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(User)

		err := BodyParser(w, r, user)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		err = s.storage.CreateUser(user)
		if err != nil {
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

func (s *Service) getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, userID)
		if userID == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		requestedUserID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		user, err := s.storage.GetUser(requestedUserID)
		if err != nil {
			code := http.StatusNotFound
			http.Error(w, http.StatusText(code), code)
			return
		}

		var resp []byte
		resp, err = json.Marshal(user)
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) deleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, userID)
		if userID == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		requestedUserID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		err = s.storage.DeleteUser(requestedUserID)
		if err != nil {
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

func (s *Service) updateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, userID)
		if userID == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		requestedUserID, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		user := new(User)

		err = BodyParser(w, r, user)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		user.ID = requestedUserID

		err = s.storage.UpdateUser(user)
		if err != nil {
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

func (s *Service) metricsHandler() http.HandlerFunc {
	log.Print("GET /metrics")
	return promhttp.Handler().(http.HandlerFunc)
}

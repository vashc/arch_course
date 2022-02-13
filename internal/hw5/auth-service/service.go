package auth_service

import (
	"arch_course/internal/hw5"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/gocraft/dbr/v2"
)

func NewService(config *Config, storage *Storage) *Service {
	return &Service{
		config:  config,
		storage: storage,
		Mux:     chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	s.Post("/register", s.registerUserHandler())

	s.Post("/login", s.loginUserHandler())

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) registerUserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(User)

		err := hw5.BodyParser(w, r, user)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		if err = s.storage.RegisterUser(user); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(RegisterResponse{ID: user.ID})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) loginUserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(LoginRequest)

		err := hw5.BodyParser(w, r, req)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Check if we have such a user in our DB
		var user *User
		user, err = s.storage.GetUserByUsername(req.Username, req.Password)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, dbr.ErrNotFound) {
				code = http.StatusNotFound
			}

			http.Error(w, http.StatusText(code), code)
			return
		}

		token := jwtauth.New(hw5.AuthTokenAlgo, []byte(s.config.JWTSecret), nil)

		// TODO: Expiration
		var tokenString string
		_, tokenString, err = token.Encode(AuthToken{hw5.AuthTokenUserID: user.ID})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		w.Header().Set(hw5.HeaderAuth, fmt.Sprintf("%s%s", hw5.HeaderBearer, tokenString))

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

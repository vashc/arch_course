package auth

import (
	"arch_course/internal/prj"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"strconv"
	"time"

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
	authUserRoute := fmt.Sprintf("/user/{%s}", prj.RequestParamUserID)

	s.Use(middleware.Timeout(5 * time.Second))

	s.Route(authUserRoute, func(router chi.Router) {
		router.Get("/", s.getUserDataHandler())
	})

	s.Post("/register", s.registerUserHandler())

	s.Post("/login", s.loginUserHandler())

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) registerUserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(prj.User)

		err := prj.BodyParser(w, r, user)
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

		err := prj.BodyParser(w, r, req)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Check if we have such a user in our DB
		var user *prj.User
		user, err = s.storage.GetUserByUsername(req.Username, req.Password)
		if err != nil {
			code := http.StatusInternalServerError
			if errors.Is(err, dbr.ErrNotFound) {
				code = http.StatusNotFound
			}

			http.Error(w, http.StatusText(code), code)
			return
		}

		token := jwtauth.New(prj.AuthTokenAlgo, []byte(s.config.JWTSecret), nil)

		// TODO: Expiration
		var tokenString string
		_, tokenString, err = token.Encode(AuthToken{prj.RequestParamUserID: user.ID})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		w.Header().Set(prj.HeaderAuth, fmt.Sprintf("%s%s", prj.HeaderBearer, tokenString))

		resp, err := json.Marshal(Response{Status: http.StatusText(http.StatusOK)})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) getUserDataHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		userID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamUserID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		var user *prj.User
		user, err = s.storage.GetUserByID(userID)
		if err != nil {
			log.Printf("storage.GetUserByID: %s\n", err.Error())
			code := http.StatusInternalServerError
			if errors.Is(err, dbr.ErrNotFound) {
				code = http.StatusNotFound
			}
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

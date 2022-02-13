package user_service

import (
	"arch_course/internal/hw5"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewService(config *Config, storage *Storage) *Service {
	return &Service{
		config:  config,
		storage: storage,
		Mux:     chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	s.Route("/user/{user_id}", func(router chi.Router) {
		router.Use(MiddlewareUserAuth)
		router.Use(MiddlewareUserCtx(s.config))
		router.Use(MiddlewareUserPermission)
		router.Get("/", s.getUserHandler())
		router.Patch("/", s.updateUserHandler())
		router.Delete("/", s.deleteUserHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) getUserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx := r.Context()

		userID, ok := ctx.Value(hw5.AuthTokenUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		var user *User
		user, err = s.storage.GetUserByID(userID)
		if err != nil {
			code := http.StatusInternalServerError
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

func (s *Service) updateUserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(hw5.AuthTokenUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		user := new(User)

		err := hw5.BodyParser(w, r, user)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		user.ID = userID

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

func (s *Service) deleteUserHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, ok := ctx.Value(hw5.AuthTokenUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		err := s.storage.DeleteUser(userID)
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

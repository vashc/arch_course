package order_service

import (
	"arch_course/internal/hw6"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewService(
	config *hw6.Config,
	storage *Storage,
	client *hw6.Client,
) (*Service, error) {
	// Orders queue instantiating
	if err := client.CreateQueue(hw6.QueueOrders); err != nil {
		return nil, err
	}

	return &Service{
		Config:  config,
		Storage: storage,
		Client:  client,
		Mux:     chi.NewRouter(),
	}, nil
}

func (s *Service) InstantiateRoutes() {
	s.Route("/order", func(router chi.Router) {
		router.Post("/", s.createOrderHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) createOrderHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		order := new(Order)

		err := hw6.BodyParser(w, r, order)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Create order record
		if err = s.Storage.CreateOrder(order); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Publish order for processing
		if err = s.Client.Publish(hw6.QueueOrders, order); err != nil {
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

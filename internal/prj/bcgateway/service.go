package bcgateway

import (
	"arch_course/internal/prj"
	"fmt"
	"log"

	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewService(config *Config, storage *Storage, client *prj.RabbitClient) *Service {
	return &Service{
		config:  config,
		storage: storage,
		client:  client,
		Mux:     chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	sendNotificationRoute := fmt.Sprintf("/withdraw/{%s}", prj.RequestParamUserID)

	s.Route(sendNotificationRoute, func(router chi.Router) {
		router.Post("/", s.withdrawHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) withdrawHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		order := new(prj.BcgatewayOrder)

		err := prj.BodyParser(w, r, &order)
		if err != nil {
			log.Printf("bcgateway.BodyParser error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		var code int
		code, err = s.createOrder(order)
		if err != nil {
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(
			Response{
				OrderUUID: order.UUID,
				Status:    http.StatusText(http.StatusOK),
			},
		)
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

func (s *Service) createOrder(order *prj.BcgatewayOrder) (code int, err error) {
	order.Status = StatusNew

	// Create order record in the DB
	order.ID, err = s.storage.CreateOrder(order)
	if err != nil {
		log.Printf("bcgateway.CreateOrder error: %s\n", err.Error())
		code = http.StatusInternalServerError
		return
	}

	// Create exchange message in a queue
	err = s.client.Publish(prj.QueueBcgatewayOrders, order)
	if err != nil {
		log.Printf("client.Publish: %s\n", err.Error())
		code = http.StatusInternalServerError
		return
	}

	return http.StatusOK, nil
}

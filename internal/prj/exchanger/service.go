package exchanger

import (
	"arch_course/internal/prj"
	"fmt"
	"log"
	"strconv"

	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func NewService(config *Config, storage *Storage, rabbitClient *prj.RabbitClient) *Service {
	return &Service{
		config:       config,
		storage:      storage,
		rabbitClient: rabbitClient,
		Mux:          chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	sellRoute := fmt.Sprintf("/sell/{%s}", prj.RequestParamOrderID)
	buyRoute := fmt.Sprintf("/buy/{%s}", prj.RequestParamOrderID)

	s.Post(sellRoute, s.sellCryptoHandler())
	s.Post(buyRoute, s.buyCryptoHandler())

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) sellCryptoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamOrderID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt error: %s", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		order := new(prj.ExchangeOrder)

		err = prj.BodyParser(w, r, &order)
		if err != nil {
			log.Printf("Exchanger.BodyParser error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		order.OrderID = orderID

		var code int
		code, err = s.createOrder(order, TypeSellOrder)
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

func (s *Service) buyCryptoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamOrderID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt error: %s", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		order := new(prj.ExchangeOrder)

		err = prj.BodyParser(w, r, &order)
		if err != nil {
			log.Printf("Exchanger.BodyParser error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		order.OrderID = orderID

		var code int
		code, err = s.createOrder(order, TypeBuyOrder)
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

func (s *Service) createOrder(order *prj.ExchangeOrder, orderType string) (code int, err error) {
	order.Status = StatusNew
	order.Type = orderType

	// Create order record in the DB
	order.ID, err = s.storage.CreateOrder(order)
	if err != nil {
		log.Printf("Exchanger.CreateOrder error: %s\n", err.Error())
		code = http.StatusInternalServerError
		return
	}

	// Create exchange message in a queue
	err = s.rabbitClient.Publish(prj.QueueExchangeOrders, order)
	if err != nil {
		log.Printf("client.Publish: %s\n", err.Error())
		code = http.StatusInternalServerError
		return
	}

	return http.StatusOK, nil
}

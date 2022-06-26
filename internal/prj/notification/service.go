package notification

import (
	"arch_course/internal/prj"
	"log"

	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	sendNotificationRoute := fmt.Sprintf("/send/{%s}", prj.RequestParamOrderID)

	s.Route(sendNotificationRoute, func(router chi.Router) {
		router.Post("/", s.sendNotificationHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) sendNotificationHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamOrderID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		notification := new(prj.Notification)
		if err = prj.BodyParser(w, r, notification); err != nil {
			log.Printf("prj.BodyParser: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		notification.Status = StatusNew
		notification.OrderID = orderID

		notification.ID, err = s.storage.CreateNotification(notification)
		if err != nil {
			log.Printf("storage.CreateWallet: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
		}

		// Create notification message in a queue
		err = s.client.Publish(prj.QueueNotifications, notification)
		if err != nil {
			log.Printf("client.Publish: %s\n", err.Error())
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

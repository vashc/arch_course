package wallet

import (
	"arch_course/internal/prj"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func NewService(config *Config, storage *Storage, rabbitClient *prj.RabbitClient) *Service {
	return &Service{
		config:  config,
		storage: storage,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		rabbitClient: rabbitClient,
		Mux:          chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	walletUserRoute := fmt.Sprintf("/wallet/{%s}", prj.RequestParamUserID)
	walletSellRoute := fmt.Sprintf("/sell/{%s}", prj.RequestParamUserID)
	walletBuyRoute := fmt.Sprintf("/buy/{%s}", prj.RequestParamUserID)

	// Authentication and context middlewares
	s.Use(middleware.Timeout(5 * time.Second))

	// Wallet routes
	s.Route(walletUserRoute, func(router chi.Router) {
		router.Use(MiddlewareUserAuth)
		router.Use(MiddlewareUserCtx(s.config))
		router.Use(MiddlewareUserPermission)

		router.Get("/", s.getWalletHandler())
		router.Post("/", s.createWalletHandler())
		router.Patch("/", s.depositWalletHandler())
	})

	s.Route(walletSellRoute, func(router chi.Router) {
		router.Use(MiddlewareUserAuth)
		router.Use(MiddlewareUserCtx(s.config))
		router.Use(MiddlewareUserPermission)

		router.Post("/", s.sellCryptoHandler())
	})

	s.Route(walletBuyRoute, func(router chi.Router) {
		router.Use(MiddlewareUserAuth)
		router.Use(MiddlewareUserCtx(s.config))
		router.Use(MiddlewareUserPermission)

		router.Post("/", s.buyCryptoHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) getWalletHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx := r.Context()

		userID, ok := ctx.Value(prj.RequestParamUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Make request to balance service
		balanceURL := fmt.Sprintf(
			"http://%s:%d/wallet/%d",
			s.config.BalanceHost,
			s.config.BalancePort,
			userID,
		)
		req, err := http.NewRequest(http.MethodGet, balanceURL, nil)
		if err != nil {
			log.Printf("http.NewRequest error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		balanceResp, err := s.client.Do(req)
		if err != nil {
			log.Printf("client.Do error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		defer balanceResp.Body.Close()

		if balanceResp.StatusCode != http.StatusOK {
			log.Printf("balanceResp.StatusCode: %d\n", balanceResp.StatusCode)
			code := balanceResp.StatusCode
			http.Error(w, http.StatusText(code), code)
			return
		}

		body, err := ioutil.ReadAll(balanceResp.Body)
		if err != nil {
			log.Printf("ioutil.ReadAll: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		// We can just proxy balance response
		_, _ = w.Write(body)
	}
}

func (s *Service) createWalletHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx := r.Context()

		userID, ok := ctx.Value(prj.RequestParamUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Make request to balance service
		balanceURL := fmt.Sprintf(
			"http://%s:%d/wallet/%d",
			s.config.BalanceHost,
			s.config.BalancePort,
			userID,
		)
		req, err := http.NewRequest(http.MethodPost, balanceURL, nil)
		if err != nil {
			log.Printf("http.NewRequest error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		balanceResp, err := s.client.Do(req)
		if err != nil {
			log.Printf("client.Do error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		defer balanceResp.Body.Close()

		if balanceResp.StatusCode != http.StatusOK {
			log.Printf("balanceResp.StatusCode: %d\n", balanceResp.StatusCode)
			code := balanceResp.StatusCode
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(Response{Status: http.StatusText(http.StatusOK)})
		if err != nil {
			log.Printf("ioutil.ReadAll: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) depositWalletHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx := r.Context()

		userID, ok := ctx.Value(prj.RequestParamUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		wallet := new(prj.Wallet)

		err = prj.BodyParser(w, r, wallet)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		var code int
		code, err = updateBalance(s, userID, 0, wallet.FiatAmount)
		if err != nil {
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

func (s *Service) sellCryptoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		ctx := r.Context()

		userID, ok := ctx.Value(prj.RequestParamUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		order := new(Order)
		err = prj.BodyParser(w, r, order)
		if err != nil {
			log.Printf("order BodyParser error: %s", err.Error())
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Fill in all the necessary fields
		order.UserID = userID
		order.Type = TypeSellOrder
		order.FiatAmount = order.CryptoAmount * CryptoToFiatRatio
		order.Status = StatusPending

		// Create sell order record
		order.ID, err = s.storage.CreateOrder(order)
		if err != nil {
			log.Printf("storage.CreateOrder error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		defer func() {
			if err != nil {
				err = s.storage.UpdateOrderStatus(order.ID, StatusFailed)
				if err != nil {
					log.Printf("storage.UpdateOrderStatus to StatusFailed: %s\n", err.Error())
				}
			}
		}()

		// Hold fiat
		code, err = holdCrypto(s, order, false)
		if err != nil {
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Create order processing record
		err = s.storage.CreateOrderProcessing(&OrderProcessing{
			OrderID:     order.ID,
			StepsNumber: 1,
		})
		if err != nil {
			log.Printf("storage.CreateOrderProcessing error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Create exchanger sell request
		code, err = createExchangerOrder(s, order, false)
		if err != nil {
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(
			Response{
				OrderID: order.ID,
				Status:  http.StatusText(http.StatusOK),
			},
		)
		if err != nil {
			log.Printf("json.Marshal: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) buyCryptoHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			code int
			err  error
		)

		ctx := r.Context()

		userID, ok := ctx.Value(prj.RequestParamUserID).(int64)
		if !ok {
			code := http.StatusUnauthorized
			http.Error(w, http.StatusText(code), code)
			return
		}

		order := new(Order)
		err = prj.BodyParser(w, r, order)
		if err != nil {
			log.Printf("order BodyParser error: %s", err.Error())
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Fill in all the necessary fields
		order.UserID = userID
		order.Type = TypeBuyOrder
		order.CryptoAmount = order.FiatAmount / CryptoToFiatRatio
		order.Status = StatusPending

		// Create sell order record
		order.ID, err = s.storage.CreateOrder(order)
		if err != nil {
			log.Printf("storage.CreateOrder error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}
		defer func() {
			if err != nil {
				err = s.storage.UpdateOrderStatus(order.ID, StatusFailed)
				if err != nil {
					log.Printf("storage.UpdateOrderStatus to StatusFailed: %s\n", err.Error())
				}
			}
		}()

		// Hold fiat
		code, err = holdFiat(s, order, false)
		if err != nil {
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Create order processing record
		err = s.storage.CreateOrderProcessing(&OrderProcessing{
			OrderID:     order.ID,
			StepsNumber: 2,
		})
		if err != nil {
			log.Printf("storage.CreateOrderProcessing error: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Create exchanger sell request
		code, err = createExchangerOrder(s, order, false)
		if err != nil {
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Create bcgateway withdraw request
		code, err = createBcgatewayOrder(s, order, false)
		if err != nil {
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(
			Response{
				OrderID: order.ID,
				Status:  http.StatusText(http.StatusOK),
			},
		)
		if err != nil {
			log.Printf("json.Marshal: %s\n", err.Error())
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

func holdFiat(s ServiceProvider, order *Order, compensate bool) (code int, err error) {
	// Decrease user fiat balance
	code, err = updateBalance(s, order.UserID, 0, -prj.Btof(compensate)*order.FiatAmount)
	if err != nil || code != http.StatusOK {
		return code, err
	}

	// Increase hot wallet fiat balance
	code, err = updateBalance(s, prj.HotWalletUserID, 0, prj.Btof(compensate)*order.FiatAmount)
	if err != nil || code != http.StatusOK {
		return code, err
	}

	return http.StatusOK, nil
}

func holdCrypto(s ServiceProvider, order *Order, compensate bool) (code int, err error) {
	// Decrease user crypto balance
	code, err = updateBalance(s, order.UserID, -prj.Btof(compensate)*order.CryptoAmount, 0)
	if err != nil || code != http.StatusOK {
		return code, err
	}

	// Increase hot wallet crypto balance
	code, err = updateBalance(s, prj.HotWalletUserID, prj.Btof(compensate)*order.CryptoAmount, 0)
	if err != nil || code != http.StatusOK {
		return code, err
	}

	return http.StatusOK, nil
}

func createExchangerOrder(s ServiceProvider, order *Order, compensate bool) (code int, err error) {
	exchangerURL := fmt.Sprintf(
		"%s/%s/%d",
		s.ExchangerURI(),
		order.Type,
		order.ID,
	)

	exchangerResp, err := prj.DoRequest(
		s.HttpClient(),
		exchangerURL,
		http.MethodPost,
		&prj.ExchangeOrder{
			AcquirerUserID: order.UserID,
			OrderID:        order.ID,
			CryptoAmount:   order.CryptoAmount,
			FiatAmount:     order.FiatAmount,
			Compensate:     compensate,
		},
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		log.Printf("exchanger DoRequest error: %s\n", err.Error())
		code = http.StatusInternalServerError
		return
	}

	if exchangerResp.StatusCode != http.StatusOK {
		log.Printf("exchangerResp.StatusCode: %d\n", exchangerResp.StatusCode)
		code = exchangerResp.StatusCode
		statusText := http.StatusText(code)
		err = errors.New(statusText)
		return
	}

	return http.StatusOK, nil
}

func createBcgatewayOrder(s ServiceProvider, order *Order, compensate bool) (code int, err error) {
	bcgatewayURL := fmt.Sprintf(
		"%s/withdraw/%d",
		s.BcgatewayURI(),
		order.UserID,
	)

	bcgatewayResp, err := prj.DoRequest(
		s.HttpClient(),
		bcgatewayURL,
		http.MethodPost,
		&prj.BcgatewayOrder{
			AcquirerUserID: order.UserID,
			OrderID:        order.ID,
			CryptoAmount:   order.CryptoAmount,
			Compensate:     compensate,
		},
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		log.Printf("exchanger DoRequest error: %s\n", err.Error())
		code = http.StatusInternalServerError
		return
	}

	if bcgatewayResp.StatusCode != http.StatusOK {
		log.Printf("bcgatewayResp.StatusCode: %d\n", bcgatewayResp.StatusCode)
		code = bcgatewayResp.StatusCode
		statusText := http.StatusText(code)
		err = errors.New(statusText)
		return
	}

	return http.StatusOK, nil
}

func updateBalance(s ServiceProvider, userID int64, cryptoAmount, fiatAmount float64) (code int, err error) {
	// Decrease user fiat balance
	balanceURL := fmt.Sprintf(
		"%s/wallet/%d",
		s.BalanceURI(),
		userID,
	)

	balanceResp, err := prj.DoRequest(
		s.HttpClient(),
		balanceURL,
		http.MethodPatch,
		&prj.Wallet{
			CryptoAmount: cryptoAmount,
			FiatAmount:   fiatAmount,
		},
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		log.Printf("balance DoRequest error: %s\n", err.Error())
		code = http.StatusInternalServerError
		return code, errors.New(http.StatusText(code))
	}

	if balanceResp.StatusCode != http.StatusOK {
		log.Printf("balanceResp.StatusCode: %d\n", balanceResp.StatusCode)
		code = balanceResp.StatusCode
		return code, errors.New(http.StatusText(code))
	}

	return http.StatusOK, nil
}

func (s *Service) HttpClient() *http.Client {
	return s.client
}

func (s *Service) BalanceURI() string {
	return fmt.Sprintf("http://%s:%d", s.config.BalanceHost, s.config.BalancePort)
}

func (s *Service) BcgatewayURI() string {
	return fmt.Sprintf("http://%s:%d", s.config.BcgatewayHost, s.config.BcgatewayPort)
}

func (s *Service) ExchangerURI() string {
	return fmt.Sprintf("http://%s:%d", s.config.ExchangerHost, s.config.ExchangerPort)
}

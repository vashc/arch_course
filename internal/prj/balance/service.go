package balance

import (
	"arch_course/internal/prj"
	"log"

	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	walletUserRoute := fmt.Sprintf("/wallet/{%s}", prj.RequestParamUserID)

	s.Route(walletUserRoute, func(router chi.Router) {
		router.Get("/", s.getWalletHandler())
		router.Post("/", s.createWalletHandler())
		router.Patch("/", s.updateWalletHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) getWalletHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		userID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamUserID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		var wallet *prj.Wallet
		wallet, err = s.storage.GetWalletByUserID(userID)
		if err != nil {
			log.Printf("storage.GetWalletByUserID: %s\n", err.Error())
			code := http.StatusInternalServerError
			if errors.Is(err, dbr.ErrNotFound) {
				code = http.StatusNotFound
			}
			http.Error(w, http.StatusText(code), code)
			return
		}

		var resp []byte
		resp, err = json.Marshal(wallet)
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) createWalletHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamUserID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		var wallet *prj.Wallet
		wallet, err = s.storage.GetWalletByUserID(userID)
		if err != nil {
			// Wallet has not been created yet
			if errors.Is(err, dbr.ErrNotFound) {
				wallet = &prj.Wallet{
					UserID:       userID,
					CryptoAmount: 0,
					FiatAmount:   0,
				}

				err = s.storage.CreateWallet(wallet)
				if err != nil {
					log.Printf("storage.CreateWallet: %s\n", err.Error())
					code := http.StatusInternalServerError
					http.Error(w, http.StatusText(code), code)
					return
				}

				resp, err := json.Marshal(Response{Status: http.StatusText(http.StatusOK)})
				if err != nil {
					log.Printf("json.Marshal: %s\n", err.Error())
					code := http.StatusInternalServerError
					http.Error(w, http.StatusText(code), code)
					return
				}

				_, _ = w.Write(resp)
				return
			}

			log.Printf("storage.GetWalletByUserID: %s\n", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Wallet for such user_id exists
		code := http.StatusConflict
		http.Error(w, http.StatusText(code), code)
		return
	}
}

func (s *Service) updateWalletHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(chi.URLParam(r, prj.RequestParamUserID), 10, 64)
		if err != nil {
			log.Printf("strconv.ParseInt error: %s", err.Error())
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		newWallet := new(prj.Wallet)

		err = prj.BodyParser(w, r, newWallet)
		if err != nil {
			log.Printf("prj.BodyParser error: %s", err.Error())
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		newWallet.UserID = userID

		// Negative amount validation
		wallet, err := s.storage.GetWalletByUserID(userID)
		if err != nil {
			log.Printf("storage.GetWalletByUserID error: %s", err.Error())
			code := http.StatusInternalServerError
			if errors.Is(err, dbr.ErrNotFound) {
				code = http.StatusNotFound
			}
			http.Error(w, http.StatusText(code), code)
			return
		}

		wallet.CryptoAmount = wallet.CryptoAmount + newWallet.CryptoAmount
		wallet.FiatAmount = wallet.FiatAmount + newWallet.FiatAmount
		if wallet.CryptoAmount < 0 || wallet.FiatAmount < 0 {
			log.Printf(
				"negative amont detected. New fiat: %f, new crypto: %f",
				wallet.FiatAmount,
				wallet.CryptoAmount,
			)
			code := http.StatusLocked
			http.Error(w, http.StatusText(code), code)
			return
		}

		err = s.storage.UpdateWallet(wallet)
		if err != nil {
			log.Printf("storage.UpdateWallet error: %s", err.Error())
			code := http.StatusInternalServerError
			if errors.Is(err, dbr.ErrNotFound) {
				code = http.StatusNotFound
			}
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

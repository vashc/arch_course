package billing_service

import (
	"arch_course/internal/hw6"
	"strconv"

	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewService(config *hw6.Config, storage *Storage) *Service {
	return &Service{
		Config:  config,
		Storage: storage,
		Mux:     chi.NewRouter(),
	}
}

func (s *Service) InstantiateRoutes() {
	s.Route("/account", func(router chi.Router) {
		router.Post("/", s.createAccountHandler())
		router.Post("/{account_id}/deposit", s.depositAccountHandler())
		router.Post("/{account_id}/withdraw", s.withdrawAccountHandler())
	})

	s.Get("/health", s.healthHandler())
}

func (s *Service) Start(port string) error {
	return http.ListenAndServe(port, s)
}

func (s *Service) createAccountHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		account := new(Account)

		err := hw6.BodyParser(w, r, account)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		if err = s.Storage.CreateAccount(account); err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		resp, err := json.Marshal(RegisterResponse{ID: account.ID})
		if err != nil {
			code := http.StatusInternalServerError
			http.Error(w, http.StatusText(code), code)
			return
		}

		_, _ = w.Write(resp)
	}
}

func (s *Service) depositAccountHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get and parse account ID
		accountID := chi.URLParam(r, hw6.AccountID)
		if accountID == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		requestedAccountID, err := strconv.ParseInt(accountID, 10, 64)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Get and parse deposit request
		deposit := new(BalanceRequest)

		err = hw6.BodyParser(w, r, deposit)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		account, err := s.Storage.GetAccountByID(requestedAccountID)
		if err != nil {
			code := http.StatusNotFound
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Change and update account balance
		account.Balance += deposit.Amount
		err = s.Storage.UpdateAccountBalance(account)
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

func (s *Service) withdrawAccountHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get and parse account ID
		accountID := chi.URLParam(r, hw6.AccountID)
		if accountID == "" {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		requestedAccountID, err := strconv.ParseInt(accountID, 10, 64)
		if err != nil {
			code := http.StatusBadRequest
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Get and parse deposit request
		withdraw := new(BalanceRequest)

		err = hw6.BodyParser(w, r, withdraw)
		if err != nil {
			code := http.StatusUnprocessableEntity
			http.Error(w, http.StatusText(code), code)
			return
		}

		account, err := s.Storage.GetAccountByID(requestedAccountID)
		if err != nil {
			code := http.StatusNotFound
			http.Error(w, http.StatusText(code), code)
			return
		}

		// Compare and update account balance if possible
		err = s.Storage.CompareAndUpdateAccountBalance(account, withdraw.Amount)
		if err != nil {
			code := http.StatusInternalServerError

			if errors.Is(err, hw6.ErrInsufficientFunds) {
				code = http.StatusForbidden
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

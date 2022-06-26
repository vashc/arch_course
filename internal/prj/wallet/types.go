package wallet

import "net/http"

type ServiceProvider interface {
	HttpClient() *http.Client
	BalanceURI() string
	BcgatewayURI() string
	ExchangerURI() string
}

type BalanceUpdater interface {
	updateBalance(s ServiceProvider, userID int64, cryptoAmount, fiatAmount float64) (code int, err error)
}

type Step func(s ServiceProvider, order *Order, compensate bool) (int, error)

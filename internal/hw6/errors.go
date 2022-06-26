package hw6

import "errors"

var (
	ErrUnsupportedMediaType = errors.New("Content-Type header is not application/json")
	ErrRequestBodyDecoding  = errors.New("request body contains badly formed JSON")
	ErrInsufficientFunds    = errors.New("insufficient funds")

	ErrDbrOpenConnection = errors.New("dbr failed to create connection")

	ErrAccountNotFound = errors.New("account not found")
)

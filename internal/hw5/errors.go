package hw5

import "errors"

var (
	ErrUnsupportedMediaType = errors.New("Content-Type header is not application/json")
	ErrRequestBodyDeconding = errors.New("request body contains badly formed JSON")
	ErrUnathorizedUser      = errors.New("unauthorized user")
	ErrWrongSignMethod      = errors.New("incorrect sign method")
	ErrInvalidTokenFormat   = errors.New("invalid token format")
	ErrInvalidToken         = errors.New("invalid token")
)

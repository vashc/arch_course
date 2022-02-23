package hw3

import "errors"

var (
	errDbrOpenConnection = errors.New("dbr failed to create connection")
)

var (
	ErrUnsupportedMediaType = errors.New("Content-Type header is not application/json")
	ErrRequestBodyDeconding = errors.New("request body contains badly formed JSON")
)

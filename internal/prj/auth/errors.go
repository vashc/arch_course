package auth

import "errors"

var (
	errDbrOpenConnection = errors.New("dbr failed to create connection")
)

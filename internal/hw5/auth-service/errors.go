package auth_service

import "errors"

var (
	errDbrOpenConnection = errors.New("dbr failed to create connection")
)

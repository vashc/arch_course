package hw2

import "errors"

var (
	errDbrOpenConnection = errors.New("dbr failed to create connection")
)

var (
	errEmptyUserID           = errors.New("empty user ID")
	errIncorrectUserIDFormat = errors.New("incorrect user ID format")
)

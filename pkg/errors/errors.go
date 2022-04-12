package errors

import (
	"errors"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrIncorrectPassword = errors.New("incorrect password")
var ErrBadRequest = errors.New("bad request")
var ErrTokenExpired = errors.New("token has expired")
var ErrTokenInvalid = errors.New("invalid token")

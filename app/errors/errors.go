package errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrUserNotFound = errors.New("user not found")
var ErrSessionNotFound = errors.New("session not found")
var ErrIncorrectPassword = errors.New("incorrect password")

type GRPCError struct {
	err    error
	status *status.Status
}

func Convert(err error) GRPCError {
	switch err {
	case ErrUserNotFound:
		return newGRPCError(err, codes.NotFound)
	case ErrIncorrectPassword:
		return newGRPCError(err, codes.PermissionDenied)
	case ErrSessionNotFound:
		return newGRPCError(err, codes.Unauthenticated)
	default:
		return newGRPCError(err, codes.Internal)
	}
}

func newGRPCError(err error, code codes.Code) GRPCError {
	return GRPCError{
		err:    err,
		status: status.New(code, err.Error()),
	}
}

func (e GRPCError) Error() string {
	return e.Error()
}

func (e GRPCError) GRPCStatus() *status.Status {
	return e.status
}

func (e GRPCError) Unwrap() error {
	return e.err
}

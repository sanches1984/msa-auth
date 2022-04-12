package service

import (
	"github.com/sanches1984/auth/pkg/errors"
	dberr "github.com/sanches1984/gopkg-pg-orm/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError interface {
	Error() string
	GRPCStatus() *status.Status
	Unwrap() error
}

type grpcError struct {
	err    error
	status *status.Status
}

func convert(err error) GRPCError {
	if v, ok := err.(grpcError); ok {
		return v
	}

	if v, ok := err.(dberr.Error); ok {
		if v.TypeOf(dberr.NotFound) {
			return newGRPCError(v, codes.NotFound)
		} else if v.TypeOf(dberr.BadRequest) {
			return newGRPCError(v, codes.InvalidArgument)
		} else if v.TypeOf(dberr.Conflict) {
			return newGRPCError(v, codes.AlreadyExists)
		} else {
			return newGRPCError(v, codes.Internal)
		}
	}

	switch err {
	case errors.ErrUserNotFound:
		return newGRPCError(err, codes.NotFound)
	case errors.ErrIncorrectPassword:
		return newGRPCError(err, codes.PermissionDenied)
	case errors.ErrSessionNotFound, errors.ErrTokenExpired, errors.ErrTokenInvalid:
		return newGRPCError(err, codes.Unauthenticated)
	case errors.ErrBadRequest:
		return newGRPCError(err, codes.InvalidArgument)
	default:
		return newGRPCError(err, codes.Internal)
	}
}

func newGRPCError(err error, code codes.Code) grpcError {
	return grpcError{
		err:    err,
		status: status.New(code, err.Error()),
	}
}

func (e grpcError) Error() string {
	return e.Error()
}

func (e grpcError) GRPCStatus() *status.Status {
	return e.status
}

func (e grpcError) Unwrap() error {
	return e.err
}

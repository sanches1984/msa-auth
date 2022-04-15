package service

import (
	"errors"
	"github.com/golang/mock/gomock"
	dberr "github.com/sanches1984/gopkg-pg-orm/errors"
	errs "github.com/sanches1984/msa-auth/pkg/errors"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

type ConverterSuite struct {
	suite.Suite

	ctrl *gomock.Controller
}

func (s *ConverterSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
}

func (s *ConverterSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestConverterSuite(t *testing.T) {
	suite.Run(t, new(ConverterSuite))
}

func (s *ConverterSuite) TestConvert() {
	cases := []struct {
		err  error
		code codes.Code
		msg  string
	}{
		{
			err: grpcError{
				err:    errors.New("test"),
				status: status.New(codes.Aborted, "message"),
			},
			code: codes.Aborted,
		},
		{
			err:  dberr.NewConflictError(errors.New("test")),
			code: codes.AlreadyExists,
		},
		{
			err:  dberr.NewBadRequestError(errors.New("test")),
			code: codes.InvalidArgument,
		},
		{
			err:  dberr.NewNotFoundError(errors.New("test")),
			code: codes.NotFound,
		},
		{
			err:  dberr.NewInternalError(errors.New("test")),
			code: codes.Internal,
		},
		{
			err:  errs.ErrUserNotFound,
			code: codes.NotFound,
		},
		{
			err:  errs.ErrTokenInvalid,
			code: codes.Unauthenticated,
		},
		{
			err:  errs.ErrTokenExpired,
			code: codes.Unauthenticated,
		},
		{
			err:  errs.ErrSessionNotFound,
			code: codes.Unauthenticated,
		},
		{
			err:  errs.ErrIncorrectPassword,
			code: codes.PermissionDenied,
		},
		{
			err:  errs.ErrBadRequest,
			code: codes.InvalidArgument,
		},
		{
			err:  errors.New("test"),
			code: codes.Internal,
		},
	}

	for n, c := range cases {
		grpcErr := convert(c.err)

		s.Equalf(c.code, grpcErr.GRPCStatus().Code(), "unexpected grpc code in case %d", n)
		s.Equalf(c.err.Error(), grpcErr.Error(), "unexpected error in case %d", n)
	}
}

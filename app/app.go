package app

import (
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/app/service"
	api "github.com/sanches1984/auth/proto/api"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	logger zerolog.Logger
}

func New(logger zerolog.Logger) *App {
	return &App{logger: logger}
}

func (a *App) Serve(addr string) error {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// todo db

	s := grpc.NewServer()
	api.RegisterAuthServiceServer(s, service.NewAuthService(a.logger))
	api.RegisterManageServiceServer(s, service.NewManageService(a.logger))

	a.logger.Info().Msg("auth service started")
	return s.Serve(conn)
}

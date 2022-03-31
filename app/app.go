package app

import (
	"database/sql"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/app/factory"
	"github.com/sanches1984/auth/app/service"
	api "github.com/sanches1984/auth/proto/api"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	db     *sql.DB
	logger zerolog.Logger
}

func New(db *sql.DB, logger zerolog.Logger) *App {
	return &App{db: db, logger: logger}
}

func (a *App) Serve(addr string) error {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	factoryService := factory.New(a.db)

	s := grpc.NewServer()
	api.RegisterAuthServiceServer(s, service.NewAuthService(factoryService, a.logger))
	api.RegisterManageServiceServer(s, service.NewManageService(factoryService, a.logger))

	a.logger.Debug().Msg("auth service started")
	return s.Serve(conn)
}

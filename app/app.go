package app

import (
	"database/sql"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/app/repository"
	"github.com/sanches1984/auth/app/resources"
	"github.com/sanches1984/auth/app/service"
	"github.com/sanches1984/auth/app/storage"
	"github.com/sanches1984/auth/config"
	"github.com/sanches1984/auth/pkg/jwt"
	"github.com/sanches1984/auth/pkg/redis"
	api "github.com/sanches1984/auth/proto/api"
	"google.golang.org/grpc"
	"net"
)

type App struct {
	db    *sql.DB
	redis *redis.Client

	repo    Repository
	storage Storage
	logger  zerolog.Logger
}

func New(logger zerolog.Logger) (*App, error) {
	var err error
	app := &App{logger: logger}

	app.db, err = resources.InitDatabase(logger)
	if err != nil {
		return app, fmt.Errorf("db init error: %w", err)
	}

	app.redis, err = resources.InitRedis()
	if err != nil {
		return app, fmt.Errorf("redis init error: %w", err)
	}

	jwtService := jwt.NewService(config.AccessTokenTTL(), config.RefreshTokenTTL(), config.Secrets().JwtSecret)
	app.repo = repository.New(app.db)
	app.storage = storage.New(app.redis, jwtService)

	return app, nil
}

func (a *App) Close() {
	if a.db != nil {
		defer a.db.Close()
	}
	if a.redis != nil {
		defer a.redis.Close()
	}
}

func (a *App) Serve(addr string) error {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	api.RegisterAuthServiceServer(s, service.NewAuthService(a.repo, a.storage, a.logger))
	api.RegisterManageServiceServer(s, service.NewManageService(a.repo, a.storage, a.logger))

	a.logger.Info().Str("addr", config.Addr()).Msg("listen")
	return s.Serve(conn)
}

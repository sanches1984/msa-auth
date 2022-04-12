package app

import (
	"fmt"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/config"
	"github.com/sanches1984/auth/internal/app/resources"
	"github.com/sanches1984/auth/internal/app/service"
	"github.com/sanches1984/auth/internal/pkg/metrics"
	"github.com/sanches1984/auth/internal/pkg/repository"
	"github.com/sanches1984/auth/internal/pkg/storage"
	"github.com/sanches1984/auth/pkg/jwt"
	"github.com/sanches1984/auth/pkg/redis"
	api "github.com/sanches1984/auth/proto/api"
	database "github.com/sanches1984/gopkg-pg-orm"
	dbmw "github.com/sanches1984/gopkg-pg-orm/middleware"
	"google.golang.org/grpc"
	"net"
	"time"
)

type App struct {
	grpc    *grpc.Server
	db      database.IClient
	redis   *redis.Client
	repo    *repository.Repository
	storage *storage.Storage
	metrics *metrics.Service
	logger  zerolog.Logger
}

func New(logger zerolog.Logger) (*App, error) {
	var err error
	app := &App{logger: logger}

	app.db, err = resources.InitDatabase(config.Env().MigrationsPath, logger)
	if err != nil {
		return app, fmt.Errorf("db init error: %w", err)
	}

	app.redis, err = resources.InitRedis(logger)
	if err != nil {
		return app, fmt.Errorf("redis init error: %w", err)
	}

	jwtService := jwt.NewService(config.Env().AccessTTL, config.Env().RefreshTTL, config.Env().JwtSecret)
	app.repo = repository.New()
	app.storage = storage.New(app.redis, jwtService)
	app.metrics = metrics.NewService(config.Env().MetricsHost)

	app.grpc = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmw.ChainUnaryServer(
				dbmw.NewDBServerInterceptor(app.db, database.WithLogger(logger, time.Second)),
				app.metrics.AppMetricsInterceptor(),
				app.metrics.GRPCMetricsInterceptor(),
			),
		),
	)

	api.RegisterAuthServiceServer(app.grpc, service.NewAuthService(app.repo, app.storage, app.logger))
	api.RegisterManageServiceServer(app.grpc, service.NewManageService(app.repo, app.storage, app.logger))
	app.metrics.Initialize(app.grpc)

	return app, nil
}

func (a *App) Close() {
	if a.db != nil {
		defer a.db.Close()
	}
	if a.redis != nil {
		defer a.redis.Close()
	}
	if a.grpc != nil {
		defer a.grpc.GracefulStop()
	}
}

func (a *App) Serve(addr string) error {
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		a.logger.Info().Str("host", config.Env().MetricsHost).Msg("listen metrics")
		if err := a.metrics.Listen(); err != nil {
			a.logger.Error().Err(err).Msg("metrics failed")
		}
	}()

	a.logger.Info().Str("host", config.Env().Host).Msg("listen service")
	return a.grpc.Serve(conn)
}

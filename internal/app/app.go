package app

import (
	"fmt"
	grpcmw "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/rs/zerolog"
	database "github.com/sanches1984/gopkg-pg-orm"
	dbmw "github.com/sanches1984/gopkg-pg-orm/middleware"
	"github.com/sanches1984/msa-auth/config"
	"github.com/sanches1984/msa-auth/internal/app/resources"
	"github.com/sanches1984/msa-auth/internal/app/service"
	"github.com/sanches1984/msa-auth/internal/pkg/metrics"
	"github.com/sanches1984/msa-auth/internal/pkg/repository"
	"github.com/sanches1984/msa-auth/internal/pkg/storage"
	"github.com/sanches1984/msa-auth/pkg/jwt"
	"github.com/sanches1984/msa-auth/pkg/redis"
	api "github.com/sanches1984/msa-auth/proto/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const gracefulTimeout = 2 * time.Second

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
		app.db.Close()
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

	grpc_health_v1.RegisterHealthServer(app.grpc, health.NewServer())
	api.RegisterAuthServiceServer(app.grpc, service.NewAuthService(app.repo, app.storage, app.logger))
	api.RegisterManageServiceServer(app.grpc, service.NewManageService(app.repo, app.storage, app.logger))
	app.metrics.Initialize(app.grpc)

	return app, nil
}

func (a *App) Serve() error {
	defer a.stop()

	conn, err := net.Listen("tcp", config.Env().Host)
	if err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		a.logger.Warn().Msg("termination signal received")
		a.logger.Info().Msg("stop grpc server")
		a.grpc.GracefulStop()
	}()

	go func() {
		a.logger.Info().Str("host", config.Env().MetricsHost).Msg("start metrics server")
		if err := a.metrics.Listen(); err != nil {
			a.logger.Error().Err(err).Msg("metrics failed")
		}
	}()

	a.logger.Info().Str("host", config.Env().Host).Msg("start grpc server")
	return a.grpc.Serve(conn)
}

func (a *App) stop() {
	time.Sleep(gracefulTimeout)

	if a.db != nil {
		a.logger.Info().Msg("disconnect database")
		a.db.Close()
	}
	if a.redis != nil {
		a.logger.Info().Msg("disconnect redis")
		a.redis.Close()
	}
	if a.metrics != nil {
		a.logger.Info().Msg("stop metrics server")
		a.metrics.Close()
	}
}

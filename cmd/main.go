package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sanches1984/auth/app"
	"github.com/sanches1984/auth/app/resources"
	"github.com/sanches1984/auth/config"
	"os"
)

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	if err := config.Load(); err != nil {
		logger.Fatal().Err(err).Msg("config load error")
	}

	db, err := resources.InitDatabase(logger)
	if db != nil {
		defer db.Close()
	}
	if err != nil {
		logger.Fatal().Err(err).Msg("db init error")
	}

	logger.Info().Str("addr", config.Addr()).Msg("listen")

	if err := app.New(db, logger).Serve(config.Addr()); err != nil {
		logger.Fatal().Err(err).Msg("auth service error")
	}
}

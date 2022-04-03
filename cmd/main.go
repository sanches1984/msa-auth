package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sanches1984/auth/app"
	"github.com/sanches1984/auth/config"
	"os"
)

func main() {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	if err := config.Load(); err != nil {
		logger.Fatal().Err(err).Msg("config load error")
	}

	application, err := app.New(logger)
	defer application.Close()
	if err != nil {
		logger.Fatal().Err(err).Msg("app init error")
	}

	if err := application.Serve(config.Addr()); err != nil {
		logger.Fatal().Err(err).Msg("auth service error")
	}
}

package main

import (
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/app"
	"github.com/sanches1984/auth/config"
)

const addr = "localhost:5000"

func main() {
	logger := zerolog.Logger{}.With().Str("service", "auth").Logger()

	if err := config.Load(); err != nil {
		logger.Fatal().Err(err).Msg("config load error")
	}

	if err := app.New(logger).Serve(addr); err != nil {
		logger.Fatal().Err(err).Msg("auth service error")
	}
}

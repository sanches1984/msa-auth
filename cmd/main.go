package main

import (
	"github.com/sanches1984/gopkg-logger"
	"github.com/sanches1984/msa-auth/config"
	"github.com/sanches1984/msa-auth/internal/app"
	syslog "log"
)

func main() {
	if err := config.Load(); err != nil {
		syslog.Fatalln("load config error:", err)
	}

	log.Init(config.Env().LogType, config.Env().LogLevel)
	logger := log.For(config.Env().AppName)

	application, err := app.New(logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("app init error")
	}

	if err := application.Serve(); err != nil {
		logger.Fatal().Err(err).Msg("auth service error")
	}
}

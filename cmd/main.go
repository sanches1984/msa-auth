package main

import (
	"github.com/sanches1984/auth/config"
	"github.com/sanches1984/auth/internal/app"
	"github.com/sanches1984/auth/internal/app/resources"
	syslog "log"
)

func main() {
	if err := config.Load(); err != nil {
		syslog.Fatalln("load config error:", err)
	}

	logger := resources.InitLogger()

	application, err := app.New(logger)
	if err != nil {
		if application != nil {
			application.Close()
		}
		logger.Fatal().Err(err).Msg("app init error")
	}
	defer application.Close()

	if err := application.Serve(config.Env().Host); err != nil {
		logger.Fatal().Err(err).Msg("auth service error")
	}
}

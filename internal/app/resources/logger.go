package resources

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sanches1984/auth/config"
	"os"
	"strings"
)

func InitLogger() zerolog.Logger {
	var logger zerolog.Logger
	level, err := zerolog.ParseLevel(config.Env().LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	if strings.ToLower(config.Env().LogType) == "json" {
		logger = log.Level(level).With().Timestamp().Logger()
	} else {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout}).Level(level).With().Timestamp().Logger()
	}

	return logger
}

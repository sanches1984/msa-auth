package storage

import (
	"github.com/rs/zerolog"
)

type Redis interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type JwtService interface {
	NewAccessToken(userID, sessionID int64) (string, error)
	NewRefreshToken(userID, sessionID int64) (string, error)
	ParseToken(token string) (int64, int64, error)
}

type Storage struct {
	redis  Redis
	jwt    JwtService
	logger zerolog.Logger
}

func New(redis Redis, jwt JwtService, logger zerolog.Logger) *Storage {
	return &Storage{
		redis:  redis,
		jwt:    jwt,
		logger: logger,
	}
}

// todo

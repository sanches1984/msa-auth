package storage

import (
	"github.com/sanches1984/msa-auth/pkg/jwt"
	uuid "github.com/satori/go.uuid"
)

type Redis interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

type JwtService interface {
	NewAccessToken(userID int64, sessionID uuid.UUID) (jwt.Token, error)
	NewRefreshToken(userID int64, sessionID uuid.UUID) (jwt.Token, error)
	ParseToken(token string) (int64, uuid.UUID, error)
}

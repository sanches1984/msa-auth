package storage

import (
	uuid "github.com/satori/go.uuid"
)

type Session struct {
	ID      uuid.UUID
	UserID  int64
	Access  Token
	Refresh Token
}

type Token struct {
	Value     string
	ExpiresIn int32
}

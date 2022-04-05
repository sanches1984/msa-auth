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
	ExpiresAt int32
}

type SessionData struct {
	ID     uuid.UUID
	UserID int64
	Data   []byte
}

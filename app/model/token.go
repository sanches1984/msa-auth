package model

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type RefreshToken struct {
	tableName struct{}  `sql:"refresh_tokens"`
	ID        int64     `sql:"id,pk"`
	UserID    int64     `sql:"user_id,notnull"`
	SessionID uuid.UUID `sql:"session_id,notnull"`
	Token     string    `sql:"token,notnull"`
	ExpiresIn int32     `sql:"expires_in,notnull"`
	Created   time.Time `sql:"created,notnull"`
	Updated   time.Time `sql:"updated,notnull"`
}

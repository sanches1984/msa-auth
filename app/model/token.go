package model

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"time"
)

type RefreshTokenList []*RefreshToken

type RefreshToken struct {
	tableName struct{}  `pg:"refresh_tokens"`
	ID        int64     `pg:"id,pk"`
	UserID    int64     `pg:"user_id,notnull"`
	SessionID uuid.UUID `pg:"session_id,notnull"`
	Token     string    `pg:"token,notnull"`
	ExpiresIn int32     `pg:"expires_in,notnull"`
	Created   time.Time `pg:"created,notnull"`
	Updated   time.Time `pg:"updated,notnull"`
}

type RefreshTokenFilter struct {
	UserID    int64
	SessionID uuid.UUID
}

func (t *RefreshToken) BeforeInsert(ctx context.Context) (context.Context, error) {
	t.Created = time.Now()
	t.Updated = time.Now()
	return ctx, nil
}

func (t *RefreshToken) BeforeUpdate(ctx context.Context) (context.Context, error) {
	t.Updated = time.Now()
	return ctx, nil
}

func (t RefreshToken) IsExpired() bool {
	return int32(time.Now().Unix()) > t.ExpiresIn
}

func (tl RefreshTokenList) Sessions() []string {
	sessions := make([]string, 0, len(tl))
	for _, t := range tl {
		sessions = append(sessions, t.SessionID.String())
	}
	return sessions
}

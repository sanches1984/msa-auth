package service

import (
	"context"
	"github.com/sanches1984/gopkg-pg-orm/pager"
	"github.com/sanches1984/msa-auth/internal/app/model"
	storage2 "github.com/sanches1984/msa-auth/internal/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetUsers(ctx context.Context, filter model.UserFilter, pgr pager.Pager) (model.UserList, error)
	GetUser(ctx context.Context, filter model.UserFilter) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUserPassword(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, user *model.User) error
	GetRefreshTokens(ctx context.Context, filter model.RefreshTokenFilter) (model.RefreshTokenList, error)
	GetRefreshToken(ctx context.Context, filter model.RefreshTokenFilter) (*model.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	UpdateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, filter model.RefreshTokenFilter) error
}

type Storage interface {
	DecodeToken(token string) (int64, uuid.UUID, error)
	GetSessionData(token string) ([]byte, error)
	GetSessionDataByUUID(sessionID uuid.UUID) ([]byte, error)
	CreateSession(userID int64, userData []byte) (*storage2.Session, error)
	RefreshSession(userID int64, sessionID uuid.UUID, userData []byte) (*storage2.Session, error)
	UpdateSessionData(token string, userData []byte) error
	DeleteSession(token string) error
	DeleteSessionByUUID(sessionID uuid.UUID) error
}

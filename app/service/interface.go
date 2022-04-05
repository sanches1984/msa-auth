package service

import (
	"context"
	"github.com/sanches1984/auth/app/model"
	"github.com/sanches1984/auth/app/storage"
	"github.com/sanches1984/gopkg-pg-orm/pager"
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
	DeleteRefreshToken(ctx context.Context, filter model.RefreshTokenFilter) error
}

type Storage interface {
	GetSession(token string) (*storage.SessionData, error)
	CreateSession(userID int64, userData []byte) (*storage.Session, error)
	UpdateSession(token string, userData []byte) error
	DeleteSession(token string) error
}

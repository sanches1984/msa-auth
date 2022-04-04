package service

import (
	"context"
	"github.com/sanches1984/auth/app/model"
	"github.com/sanches1984/auth/app/storage"
	uuid "github.com/satori/go.uuid"
)

type Repository interface {
	GetUserByLogin(ctx context.Context, login string) (*model.User, error)
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) (int64, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	DeleteUser(ctx context.Context, id int64) error
	GetRefreshToken(ctx context.Context, userID int64, sessionID uuid.UUID) (*model.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, userID int64, sessionID uuid.UUID) error
}

type Storage interface {
	GetUserIDByToken(token string) (int64, error)
	CreateSession(userID int64, userData []byte) (*storage.Session, error)
	DeleteSession(token string) (int64, uuid.UUID, error)
}

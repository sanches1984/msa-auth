package service

import (
	"context"
	"github.com/sanches1984/auth/app/model"
)

type Repository interface {
	GetUsers(ctx context.Context, filter model.UserFilter) ([]model.User, error)
	GetUserByLogin(ctx context.Context, login string) (model.User, error)
	CreateUser(ctx context.Context, login, password string) (int64, error)
	UpdateUserPassword(ctx context.Context, id int64, password string) error
	DeleteUser(ctx context.Context, id int64) error
}

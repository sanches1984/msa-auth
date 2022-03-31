package factory

import (
	"context"
	"database/sql"
	"github.com/sanches1984/auth/app/model"
)

type Factory struct {
	db *sql.DB
}

func New(db *sql.DB) *Factory {
	return &Factory{db: db}
}

func (f *Factory) GetUsers(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	// todo
	return nil, nil
}

func (f *Factory) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	// todo
	return model.User{}, nil
}

func (f *Factory) CreateUser(ctx context.Context, login, password string) (int64, error) {
	// todo
	return 0, nil
}

func (f *Factory) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	// todo
	return nil
}

func (f *Factory) DeleteUser(ctx context.Context, id int64) error {
	// todo
	return nil
}

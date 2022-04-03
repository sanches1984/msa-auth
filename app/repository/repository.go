package repository

import (
	"context"
	"database/sql"
	"github.com/sanches1984/auth/app/model"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (f *Repository) GetUsers(ctx context.Context, filter model.UserFilter) ([]model.User, error) {
	// todo
	return nil, nil
}

func (f *Repository) GetUserByLogin(ctx context.Context, login string) (model.User, error) {
	// todo
	return model.User{}, nil
}

func (f *Repository) CreateUser(ctx context.Context, login, password string) (int64, error) {
	// todo
	return 0, nil
}

func (f *Repository) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	// todo
	return nil
}

func (f *Repository) DeleteUser(ctx context.Context, id int64) error {
	// todo
	return nil
}

package repository

import (
	"context"
	"database/sql"
	"github.com/sanches1984/auth/app/model"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (f *Repository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE login = ?"
	row := f.db.QueryRowContext(ctx, query, login)
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (f *Repository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE id = ?"
	row := f.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&user)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (f *Repository) CreateUser(ctx context.Context, user *model.User) (int64, error) {
	query := "INSERT INTO users(login, password_hash, created, updated) VALUES (?, ?, ?, ?)"
	err := f.db.QueryRowContext(ctx, query, user.Login, user.PasswordHash, user.Created, user.Updated).Scan(&user.ID)
	return user.ID, err
}

func (f *Repository) UpdateUserPassword(ctx context.Context, id int64, password string) error {
	query := "UPDATE users SET password_hash = ? WHERE id = ?"
	_, err := f.db.ExecContext(ctx, query, password, id)
	return err
}

func (f *Repository) DeleteUser(ctx context.Context, id int64) error {
	query := "DELETE FROM refresh_tokens WHERE user_id = ?"
	if _, err := f.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	query = "UPDATE users SET deleted = ? WHERE id = ?"
	_, err := f.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (f *Repository) GetRefreshToken(ctx context.Context, userID int64, sessionID uuid.UUID) (*model.RefreshToken, error) {
	// todo
	return nil, nil
}

func (f *Repository) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	// todo
	return nil
}

func (f *Repository) DeleteRefreshToken(ctx context.Context, userID int64, sessionID uuid.UUID) error {
	query := "DELETE FROM refresh_tokens WHERE user_id = ?"
	args := []interface{}{userID}
	if sessionID != uuid.Nil {
		query += " AND session_id = ?"
		args = append(args, sessionID)
	}

	_, err := f.db.ExecContext(ctx, query, args)
	return err
}

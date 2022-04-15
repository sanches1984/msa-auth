package repository

import (
	"context"
	"github.com/sanches1984/gopkg-pg-orm/pager"
	"github.com/sanches1984/gopkg-pg-orm/repository/dao"
	"github.com/sanches1984/gopkg-pg-orm/repository/opt"
	"github.com/sanches1984/msa-auth/internal/app/model"
	uuid "github.com/satori/go.uuid"
)

type Repository struct {
	db ORM
}

func New() *Repository {
	return &Repository{db: dao.New()}
}

func (r *Repository) GetUsers(ctx context.Context, filter model.UserFilter, pgr pager.Pager) (model.UserList, error) {
	var users []*model.User
	opts := opt.List()
	if filter.ID != 0 {
		opts = append(opts, opt.Eq("id", filter.ID))
	}
	if filter.Login != "" {
		opts = append(opts, opt.Eq("login", filter.Login))
	}
	if !filter.ShowDeleted {
		opts = append(opts, opt.IsNull("deleted"))
	}
	if pgr != nil {
		opts = append(opts, opt.Paging(pgr.GetPage(), pgr.GetPageSize()))
	}

	opts = append(opts, filter.Order.GetOptFn())

	err := r.db.FindList(ctx, &users, opts)
	return users, err
}

func (r *Repository) GetUser(ctx context.Context, filter model.UserFilter) (*model.User, error) {
	users, err := r.GetUsers(ctx, filter, nil)
	if err != nil {
		return nil, err
	} else if len(users) != 1 {
		return nil, nil
	}

	return users[0], nil
}

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.Insert(ctx, user)
}

func (r *Repository) UpdateUserPassword(ctx context.Context, user *model.User) error {
	return r.db.Update(ctx, user, "password_hash")
}

func (r *Repository) DeleteUser(ctx context.Context, user *model.User) error {
	opts := opt.List(opt.Eq("user_id", user.ID))
	if err := r.db.HardDeleteWhere(ctx, &model.RefreshToken{}, opts); err != nil {
		return err
	}

	return r.db.SoftDelete(ctx, user)
}

func (r *Repository) GetRefreshTokens(ctx context.Context, filter model.RefreshTokenFilter) (model.RefreshTokenList, error) {
	var tokens []*model.RefreshToken
	opts := opt.List()
	if filter.UserID != 0 {
		opts = append(opts, opt.Eq("user_id", filter.UserID))
	}
	if filter.SessionID != uuid.Nil {
		opts = append(opts, opt.Eq("session_id", filter.SessionID))
	}

	err := r.db.FindList(ctx, &tokens, opts)
	return tokens, err
}

func (r *Repository) GetRefreshToken(ctx context.Context, filter model.RefreshTokenFilter) (*model.RefreshToken, error) {
	tokens, err := r.GetRefreshTokens(ctx, filter)
	if err != nil {
		return nil, err
	} else if len(tokens) != 1 {
		return nil, nil
	}

	return tokens[0], nil
}

func (r *Repository) CreateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return r.db.Insert(ctx, token)
}

func (r *Repository) UpdateRefreshToken(ctx context.Context, token *model.RefreshToken) error {
	return r.db.Update(ctx, token, "token", "expires_in")
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, filter model.RefreshTokenFilter) error {
	opts := opt.List()
	opts = append(opts, opt.Eq("user_id", filter.UserID))
	if filter.SessionID != uuid.Nil {
		opts = append(opts, opt.Eq("session_id", filter.SessionID))
	}

	return r.db.HardDeleteWhere(ctx, &model.RefreshToken{}, opts)
}

package service

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/app/errors"
	"github.com/sanches1984/auth/app/model"
	api "github.com/sanches1984/auth/proto/api"
	"github.com/sanches1984/gopkg-pg-orm/pager"
	"time"
)

type ManageService struct {
	api.ManageServiceServer

	repo    Repository
	storage Storage
	logger  zerolog.Logger
}

func NewManageService(repo Repository, storage Storage, logger zerolog.Logger) *ManageService {
	return &ManageService{
		repo:    repo,
		storage: storage,
		logger:  logger,
	}
}

func (s *ManageService) CreateUser(ctx context.Context, r *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	user := &model.User{Login: r.GetLogin()}
	if err := user.SetHashByPassword(r.GetPassword()); err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't set password hash")
		return nil, errors.Convert(err)
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't create user")
		return nil, errors.Convert(err)
	}

	return &api.CreateUserResponse{Id: user.ID}, nil
}

func (s *ManageService) DeleteUser(ctx context.Context, r *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	user, err := s.repo.GetUser(ctx, model.UserFilter{ID: r.GetId()})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", r.GetId()).Msg("can't get user by id")
		return nil, errors.Convert(err)
	} else if user == nil {
		s.logger.Info().Int64("id", r.GetId()).Msg("user not found")
		return nil, errors.Convert(errors.ErrUserNotFound)
	}

	tokens, err := s.repo.GetRefreshTokens(ctx, model.RefreshTokenFilter{UserID: r.GetId()})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", r.GetId()).Msg("can't get refresh token list")
		return nil, errors.Convert(err)
	}

	for _, token := range tokens {
		s.logger.Debug().Msg(token.SessionID.String())
		// todo remove sessions
	}

	return &api.DeleteUserResponse{SessionId: tokens.Sessions()}, nil
}

func (s *ManageService) GetUserList(ctx context.Context, r *api.GetUserListRequest) (*api.GetUserListResponse, error) {
	filter := model.UserFilter{
		ID:    r.Id,
		Login: r.Login,
		Order: model.UserOrder(int(r.Order)),
	}
	pgr := pager.NewPagerWithPageSize(r.Page, r.PageSize)

	users, err := s.repo.GetUsers(ctx, filter, pgr)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get user list")
		return nil, errors.Convert(err)
	}

	userList := make([]*api.User, 0, len(users))
	for _, u := range users {
		userList = append(userList, &api.User{
			Id:      u.ID,
			Login:   u.Login,
			Created: u.Created.Format(time.RFC3339),
			Updated: u.Updated.Format(time.RFC3339),
			Deleted: u.Deleted.Format(time.RFC3339),
		})
	}

	return &api.GetUserListResponse{Users: userList}, nil
}

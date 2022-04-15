package service

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sanches1984/gopkg-pg-orm/pager"
	"github.com/sanches1984/msa-auth/internal/app/model"
	"github.com/sanches1984/msa-auth/pkg/errors"
	api "github.com/sanches1984/msa-auth/proto/api"
	"golang.org/x/sync/errgroup"
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
	if r.GetLogin() == "" || r.GetPassword() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	user := &model.User{Login: r.GetLogin()}
	if err := user.SetHashByPassword(r.GetPassword()); err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't set password hash")
		return nil, convert(err)
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't create user")
		return nil, convert(err)
	}

	s.logger.Info().Int64("user_id", user.ID).Msg("created new user")
	return &api.CreateUserResponse{UserId: user.ID}, nil
}

func (s *ManageService) DeleteUser(ctx context.Context, r *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	if r.GetUserId() == 0 {
		return nil, convert(errors.ErrBadRequest)
	}
	user, err := s.repo.GetUser(ctx, model.UserFilter{ID: r.GetUserId()})
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", r.GetUserId()).Msg("can't get user by id")
		return nil, convert(err)
	} else if user == nil {
		s.logger.Info().Int64("user_id", r.GetUserId()).Msg("user not found")
		return nil, convert(errors.ErrUserNotFound)
	}

	tokens, err := s.repo.GetRefreshTokens(ctx, model.RefreshTokenFilter{UserID: r.GetUserId()})
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", r.GetUserId()).Msg("can't get refresh token list")
		return nil, convert(err)
	}

	eg := errgroup.Group{}
	for _, token := range tokens {
		sessionID := token.SessionID
		eg.Go(func() error {
			return s.storage.DeleteSessionByUUID(sessionID)
		})
	}
	if err := eg.Wait(); err != nil {
		s.logger.Error().Err(err).Int64("user_id", r.GetUserId()).Msg("can't delete session")
		return nil, convert(err)
	}

	if err := s.repo.DeleteRefreshToken(ctx, model.RefreshTokenFilter{UserID: r.GetUserId()}); err != nil {
		s.logger.Error().Err(err).Int64("user_id", r.GetUserId()).Msg("can't delete refresh token")
		return nil, convert(err)
	}
	if err := s.repo.DeleteUser(ctx, user); err != nil {
		s.logger.Error().Err(err).Int64("user_id", r.GetUserId()).Msg("can't delete user")
		return nil, convert(err)
	}

	s.logger.Info().Int64("user_id", user.ID).Msg("deleted user")
	return &api.DeleteUserResponse{SessionId: tokens.Sessions()}, nil
}

func (s *ManageService) GetUsers(ctx context.Context, r *api.GetUsersRequest) (*api.GetUsersResponse, error) {
	filter := model.UserFilter{
		ID:    r.GetUserId(),
		Login: r.GetLogin(),
		Order: model.UserOrder(int(r.GetOrder())),
	}
	pgr := pager.NewPagerWithPageSize(r.GetPage(), r.GetPageSize())

	users, err := s.repo.GetUsers(ctx, filter, pgr)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get user list")
		return nil, convert(err)
	}

	userList := make([]*api.User, 0, len(users))
	for _, u := range users {
		user := &api.User{
			Id:      u.ID,
			Login:   u.Login,
			Created: u.Created.Format(time.RFC3339),
			Updated: u.Updated.Format(time.RFC3339),
		}
		if u.Deleted != nil {
			user.Deleted = u.Deleted.Format(time.RFC3339)
		}
		userList = append(userList, user)
	}

	s.logger.Info().Int("count", len(userList)).Msg("get user list")
	return &api.GetUsersResponse{Users: userList}, nil
}

package service

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/app/errors"
	"github.com/sanches1984/auth/app/model"
	"github.com/sanches1984/auth/pkg/jwt"
	api "github.com/sanches1984/auth/proto/api"
)

type AuthService struct {
	api.AuthServiceServer

	repo    Repository
	storage Storage
	logger  zerolog.Logger
}

func NewAuthService(repo Repository, storage Storage, logger zerolog.Logger) *AuthService {
	return &AuthService{
		repo:    repo,
		storage: storage,
		logger:  logger,
	}
}

func (s *AuthService) Login(ctx context.Context, r *api.LoginRequest) (*api.TokenResponse, error) {
	user, err := s.repo.GetUser(ctx, model.UserFilter{Login: r.GetLogin()})
	if err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't get user by login")
		return nil, errors.Convert(err)
	} else if user == nil {
		s.logger.Info().Str("login", r.GetLogin()).Msg("user not found")
		return nil, errors.Convert(errors.ErrUserNotFound)
	}

	if !user.IsPasswordCorrect(r.GetPassword()) {
		return nil, errors.Convert(errors.ErrIncorrectPassword)
	}

	session, err := s.storage.CreateSession(user.ID, r.GetData())
	if err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't create session")
		return nil, errors.Convert(err)
	}

	if err := s.repo.CreateRefreshToken(ctx, &model.RefreshToken{
		UserID:    session.UserID,
		SessionID: session.ID,
		Token:     session.Refresh.Value,
		ExpiresIn: session.Refresh.ExpiresAt,
	}); err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't create refresh token")
		return nil, errors.Convert(err)
	}

	return &api.TokenResponse{
		SessionId: session.ID.String(),
		Access: &api.Token{
			Token:     session.Access.Value,
			ExpiresIn: session.Access.ExpiresAt,
		},
		Refresh: &api.Token{
			Token:     session.Refresh.Value,
			ExpiresIn: session.Refresh.ExpiresAt,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, r *api.LogoutRequest) (*api.LogoutResponse, error) {
	session, err := s.storage.GetSession(r.GetToken())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get session")
		return nil, errors.Convert(err)
	} else if session == nil {
		s.logger.Warn().Msg("session not found")
		return &api.LogoutResponse{SessionId: ""}, nil
	}

	if err := s.storage.DeleteSession(r.GetToken()); err != nil {
		s.logger.Error().Err(err).Msg("can't delete session")
		return nil, errors.Convert(err)
	}

	if err := s.repo.DeleteRefreshToken(ctx, model.RefreshTokenFilter{UserID: session.UserID, SessionID: session.ID}); err != nil {
		s.logger.Error().Err(err).Msg("can't delete refresh token")
		return nil, errors.Convert(err)
	}

	return &api.LogoutResponse{SessionId: session.ID.String()}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	session, err := s.storage.GetSession(r.GetToken())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get user id")
		return nil, errors.Convert(err)
	} else if session == nil {
		s.logger.Info().Msg("session not found")
		return nil, errors.Convert(errors.ErrSessionNotFound)
	}

	user, err := s.repo.GetUser(ctx, model.UserFilter{ID: session.UserID})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", session.UserID).Msg("can't get user by id")
		return nil, errors.Convert(err)
	} else if user == nil {
		s.logger.Info().Int64("id", session.UserID).Msg("user not found")
		return nil, errors.Convert(errors.ErrUserNotFound)
	}

	err = user.SetHashByPassword(r.GetNewPassword())
	if err != nil {
		s.logger.Error().Err(err).Int64("id", session.UserID).Msg("can't set password hash")
		return nil, errors.Convert(err)
	}

	err = s.repo.UpdateUserPassword(ctx, user)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", session.UserID).Msg("can't change user password")
		return nil, errors.Convert(err)
	}

	s.logger.Debug().Int64("id", session.UserID).Msg("password changed")
	return &api.ChangePasswordResponse{Changed: true}, nil
}

func (s *AuthService) GetAccessTokenByRefreshToken(ctx context.Context, r *api.GetAccessTokenByRefreshTokenRequest) (*api.TokenResponse, error) {
	// todo
	return &api.TokenResponse{
		Access: &api.Token{
			Token:     r.RefreshToken,
			ExpiresIn: 0,
		},
		Refresh: &api.Token{
			Token:     "hello",
			ExpiresIn: 0,
		},
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, r *api.ValidateTokenRequest) (*api.ValidateTokenResponse, error) {
	session, err := s.storage.GetSession(r.GetToken())
	if err != nil {
		if err == jwt.ErrInvalidToken {
			return &api.ValidateTokenResponse{Valid: false}, nil
		}
		s.logger.Error().Err(err).Msg("can't user id")
		return nil, errors.Convert(err)
	} else if session == nil {
		return &api.ValidateTokenResponse{Valid: false}, nil
	}

	user, err := s.repo.GetUser(ctx, model.UserFilter{ID: session.UserID})
	if err != nil {
		s.logger.Error().Err(err).Int64("id", session.UserID).Msg("can't get user by id")
		return nil, errors.Convert(err)
	} else if user == nil {
		return &api.ValidateTokenResponse{Valid: false}, nil
	}

	return &api.ValidateTokenResponse{
		Valid:     true,
		UserId:    session.UserID,
		SessionId: session.ID.String(),
		Data:      session.Data,
	}, nil
}

func (s *AuthService) UpdateSessionData(ctx context.Context, r *api.UpdateSessionDataRequest) (*api.UpdateSessionDataResponse, error) {
	err := s.storage.UpdateSession(r.GetToken(), r.GetData())
	if err != nil {
		return nil, errors.Convert(err)
	}
	return &api.UpdateSessionDataResponse{Updated: true}, nil
}

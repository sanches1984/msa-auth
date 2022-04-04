package service

import (
	"context"
	"github.com/rs/zerolog"
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
	user, err := s.repo.GetUserByLogin(ctx, r.Login)
	if err != nil {
		s.logger.Error().Err(err).Str("login", r.Login).Msg("can't get user by login")
		return nil, err
	} else if user == nil {
		s.logger.Info().Str("login", r.Login).Msg("user not found")
		return nil, ErrUserNotFound
	} else if user.Deleted != nil {
		s.logger.Warn().Str("login", r.Login).Msg("user is deleted")
		return nil, ErrUserIsDeleted
	}

	if !user.IsPasswordCorrect(r.Password) {
		return nil, ErrIncorrectPassword
	}

	// todo: check refresh token

	session, err := s.storage.CreateSession(user.ID, r.Data)
	if err != nil {
		s.logger.Error().Err(err).Str("login", r.Login).Msg("can't create session")
		return nil, err
	}

	// todo: write refresh token

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
	userID, sessionID, err := s.storage.DeleteSession(r.Token)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't delete session")
		return nil, err
	}

	if err := s.repo.DeleteRefreshToken(ctx, userID, sessionID); err != nil {
		s.logger.Error().Err(err).Msg("can't delete refresh token")
		return nil, err
	}

	return &api.LogoutResponse{SessionId: sessionID.String()}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	userID, err := s.storage.GetUserIDByToken(r.Token)
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get user id")
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", userID).Msg("can't get user by id")
		return nil, err
	} else if user == nil {
		s.logger.Info().Int64("id", userID).Msg("user not found")
		return nil, ErrUserNotFound
	} else if user.Deleted != nil {
		s.logger.Warn().Int64("id", userID).Msg("user is deleted")
		return nil, ErrUserIsDeleted
	}

	err = user.SetHashByPassword(r.NewPassword)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", userID).Msg("can't set password hash")
		return nil, err
	}

	err = s.repo.UpdateUserPassword(ctx, user.ID, user.PasswordHash)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", userID).Msg("can't change user password")
		return nil, err
	}

	s.logger.Debug().Int64("id", userID).Msg("password changed")
	return &api.ChangePasswordResponse{Changed: true}, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, r *api.RefreshTokensRequest) (*api.TokenResponse, error) {
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
	userID, err := s.storage.GetUserIDByToken(r.Token)
	if err != nil {
		if err == jwt.ErrInvalidToken {
			return &api.ValidateTokenResponse{Valid: false}, nil
		}
		s.logger.Error().Err(err).Msg("can't user id")
		return nil, err
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", userID).Msg("can't get user by id")
		return nil, err
	} else if user == nil || user.Deleted != nil {
		return &api.ValidateTokenResponse{Valid: false}, nil
	}

	return &api.ValidateTokenResponse{Valid: true}, nil
}

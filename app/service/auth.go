package service

import (
	"context"
	"github.com/rs/zerolog"
	api "github.com/sanches1984/auth/proto/api"
)

type AuthService struct {
	api.AuthServiceServer

	repo   Repository
	logger zerolog.Logger
}

func NewAuthService(factory Repository, logger zerolog.Logger) *AuthService {
	return &AuthService{
		repo:   factory,
		logger: logger,
	}
}

func (s *AuthService) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	// todo redis

	err := s.repo.UpdateUserPassword(ctx, 0, r.NewPassword)
	if err != nil {
		s.logger.Error().Err(err).Int64("id", 0).Msg("can't change user password")
		return nil, err
	}

	s.logger.Debug().Int64("id", 0).Msg("password changed")
	return &api.ChangePasswordResponse{Changed: true}, nil
}

func (s *AuthService) Authorize(ctx context.Context, r *api.AuthorizeRequest) (*api.TokenResponse, error) {
	// todo
	return &api.TokenResponse{
		Access: &api.Token{
			Token:     r.Login,
			ExpiresIn: 0,
		},
		Refresh: &api.Token{
			Token:     r.Password,
			ExpiresIn: 0,
		},
	}, nil
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
	// todo
	return &api.ValidateTokenResponse{Valid: true}, nil
}

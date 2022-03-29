package service

import (
	"context"
	"github.com/rs/zerolog"
	api "github.com/sanches1984/auth/proto/api"
)

type AuthService struct {
	api.AuthServiceServer

	logger zerolog.Logger
}

func NewAuthService(logger zerolog.Logger) *AuthService {
	return &AuthService{logger: logger}
}

func (s *AuthService) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	// todo
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

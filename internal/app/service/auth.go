package service

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/internal/app/model"
	"github.com/sanches1984/auth/pkg/errors"
	"github.com/sanches1984/auth/pkg/redis"
	api "github.com/sanches1984/auth/proto/api"
	"time"
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
	if r.GetLogin() == "" || r.GetPassword() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	user, err := s.repo.GetUser(ctx, model.UserFilter{Login: r.GetLogin()})
	if err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't get user by login")
		return nil, convert(err)
	} else if user == nil {
		s.logger.Info().Str("login", r.GetLogin()).Msg("user not found")
		return nil, convert(errors.ErrUserNotFound)
	}

	if !user.IsPasswordCorrect(r.GetPassword()) {
		return nil, convert(errors.ErrIncorrectPassword)
	}

	session, err := s.storage.CreateSession(user.ID, r.GetData())
	if err != nil {
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't create session")
		return nil, convert(err)
	}

	if err := s.repo.CreateRefreshToken(ctx, &model.RefreshToken{
		UserID:    session.UserID,
		SessionID: session.ID,
		Token:     session.Refresh.Value,
		ExpiresIn: session.Refresh.ExpiresIn,
	}); err != nil {
		_ = s.storage.DeleteSession(session.Access.Value)
		s.logger.Error().Err(err).Str("login", r.GetLogin()).Msg("can't create refresh token")
		return nil, convert(err)
	}

	s.logger.Info().Int64("user_id", session.UserID).Msg("login")
	return &api.TokenResponse{
		SessionId: session.ID.String(),
		Access: &api.Token{
			Token:     session.Access.Value,
			ExpiresIn: session.Access.ExpiresIn,
		},
		Refresh: &api.Token{
			Token:     session.Refresh.Value,
			ExpiresIn: session.Refresh.ExpiresIn,
		},
	}, nil
}

func (s *AuthService) Logout(ctx context.Context, r *api.LogoutRequest) (*api.LogoutResponse, error) {
	if r.GetToken() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	userID, sessionID, err := s.storage.DecodeToken(r.GetToken())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't get session")
		return nil, convert(err)
	}

	if err := s.storage.DeleteSession(r.GetToken()); err != nil {
		s.logger.Error().Err(err).Msg("can't delete session")
		return nil, convert(err)
	}

	if err := s.repo.DeleteRefreshToken(ctx, model.RefreshTokenFilter{UserID: userID, SessionID: sessionID}); err != nil {
		s.logger.Error().Err(err).Msg("can't delete refresh token")
		return nil, convert(err)
	}

	s.logger.Info().Int64("user_id", userID).Msg("logout")
	return &api.LogoutResponse{SessionId: sessionID.String()}, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, r *api.ChangePasswordRequest) (*api.ChangePasswordResponse, error) {
	if r.GetToken() == "" || r.GetNewPassword() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	userID, _, err := s.storage.DecodeToken(r.GetToken())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't decode token")
		return nil, convert(err)
	}

	_, err = s.storage.GetSessionData(r.GetToken())
	if err != nil {
		if err == redis.ErrRecordNotFound {
			s.logger.Info().Int64("user_id", userID).Msg("session not found")
			return nil, convert(errors.ErrSessionNotFound)
		}
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't get session data")
		return nil, convert(err)
	}

	user, err := s.repo.GetUser(ctx, model.UserFilter{ID: userID})
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't get user by id")
		return nil, err
	} else if user == nil {
		s.logger.Info().Int64("user_id", userID).Msg("user not found")
		return nil, convert(errors.ErrUserNotFound)
	}

	err = user.SetHashByPassword(r.GetNewPassword())
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't set password hash")
		return nil, convert(err)
	}

	err = s.repo.UpdateUserPassword(ctx, user)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't change user password")
		return nil, convert(err)
	}

	s.logger.Info().Int64("user_id", userID).Msg("password changed")
	return &api.ChangePasswordResponse{Changed: true}, nil
}

func (s *AuthService) NewAccessTokenByRefreshToken(ctx context.Context, r *api.NewAccessTokenByRefreshTokenRequest) (*api.TokenResponse, error) {
	if r.GetRefreshToken() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	userID, sessionID, err := s.storage.DecodeToken(r.GetRefreshToken())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't decode token")
		return nil, convert(err)
	}

	refreshToken, err := s.repo.GetRefreshToken(ctx, model.RefreshTokenFilter{UserID: userID, SessionID: sessionID})
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't get refresh token")
		return nil, convert(err)
	} else if refreshToken == nil {
		s.logger.Warn().Int64("user_id", userID).Msg("refresh token not found")
		return nil, convert(errors.ErrBadRequest)
	} else if refreshToken.Token != r.GetRefreshToken() {
		return nil, convert(errors.ErrTokenInvalid)
	} else if refreshToken.IsExpired() {
		s.logger.Warn().Int64("user_id", userID).Msg("refresh token has expired")
		return nil, convert(errors.ErrTokenExpired)
	}

	sessionData, err := s.storage.GetSessionDataByUUID(sessionID)
	if err != nil {
		s.logger.Warn().Err(err).Int64("user_id", userID).Msg("can't get session data")
	}

	session, err := s.storage.RefreshSession(userID, sessionID, sessionData)
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't refresh session")
		return nil, convert(err)
	}

	refreshToken.Token = session.Refresh.Value
	refreshToken.ExpiresIn = session.Refresh.ExpiresIn
	if err := s.repo.UpdateRefreshToken(ctx, refreshToken); err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't update refresh token")
		return nil, convert(err)
	}

	s.logger.Info().Int64("user_id", userID).Msg("created new access token by refresh token")
	return &api.TokenResponse{
		SessionId: sessionID.String(),
		Access: &api.Token{
			Token:     session.Access.Value,
			ExpiresIn: session.Access.ExpiresIn,
		},
		Refresh: &api.Token{
			Token:     session.Refresh.Value,
			ExpiresIn: session.Refresh.ExpiresIn,
		},
	}, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, r *api.ValidateTokenRequest) (*api.ValidateTokenResponse, error) {
	if r.GetToken() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	userID, sessionID, err := s.storage.DecodeToken(r.GetToken())
	if err != nil {
		return nil, convert(errors.ErrTokenInvalid)
	}
	sessionData, err := s.storage.GetSessionData(r.GetToken())
	if err != nil {
		s.logger.Warn().Err(err).Int64("user_id", userID).Msg("can't get session data")
		return nil, convert(errors.ErrTokenInvalid)
	}

	user, err := s.repo.GetUser(ctx, model.UserFilter{ID: userID})
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't get user by id")
		return nil, err
	} else if user == nil {
		return nil, convert(errors.ErrTokenInvalid)
	}

	s.logger.Info().Int64("user_id", userID).Msg("validate token")
	return &api.ValidateTokenResponse{
		UserId:    userID,
		SessionId: sessionID.String(),
		Data:      sessionData,
	}, nil
}

func (s *AuthService) UpdateSessionData(ctx context.Context, r *api.UpdateSessionDataRequest) (*api.UpdateSessionDataResponse, error) {
	if r.GetToken() == "" {
		return nil, convert(errors.ErrBadRequest)
	}
	err := s.storage.UpdateSessionData(r.GetToken(), r.GetData())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't update session data")
		return nil, convert(err)
	}

	s.logger.Info().Msg("update session data")
	return &api.UpdateSessionDataResponse{Updated: true}, nil
}

func (s *AuthService) GetUserSessions(ctx context.Context, r *api.GetUserSessionsRequest) (*api.GetUserSessionsResponse, error) {
	if r.GetToken() == "" {
		return nil, convert(errors.ErrBadRequest)
	}

	userID, _, err := s.storage.DecodeToken(r.GetToken())
	if err != nil {
		s.logger.Error().Err(err).Msg("can't decode token")
		return nil, convert(err)
	}

	refreshTokens, err := s.repo.GetRefreshTokens(ctx, model.RefreshTokenFilter{UserID: userID})
	if err != nil {
		s.logger.Error().Err(err).Int64("user_id", userID).Msg("can't get refresh token")
		return nil, convert(err)
	}

	sessions := make([]*api.Session, 0, len(refreshTokens))
	for _, t := range refreshTokens {
		sessions = append(sessions, &api.Session{
			Id:      t.SessionID.String(),
			Created: t.Created.Format(time.RFC3339),
		})
	}

	return &api.GetUserSessionsResponse{Sessions: sessions}, nil
}

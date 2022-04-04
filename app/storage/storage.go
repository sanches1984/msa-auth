package storage

import (
	uuid "github.com/satori/go.uuid"
)

type Storage struct {
	redis Redis
	jwt   JwtService
}

func New(redis Redis, jwt JwtService) *Storage {
	return &Storage{
		redis: redis,
		jwt:   jwt,
	}
}

func (s *Storage) GetUserIDByToken(token string) (int64, error) {
	userID, _, err := s.jwt.ParseToken(token)
	return userID, err
}

func (s *Storage) CreateSession(userID int64, userData []byte) (*Session, error) {
	sessionID := uuid.NewV4()
	access, err := s.jwt.NewAccessToken(userID, sessionID)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwt.NewRefreshToken(userID, sessionID)
	if err != nil {
		return nil, err
	}

	session := &Session{
		ID:     sessionID,
		UserID: userID,
		Access: Token{
			Value:     access.Value,
			ExpiresAt: access.ExpiresAt,
		},
		Refresh: Token{
			Value:     refresh.Value,
			ExpiresAt: refresh.ExpiresAt,
		},
	}

	err = s.redis.Set(access.Value, userData)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Storage) DeleteSession(token string) (int64, uuid.UUID, error) {
	userID, sessionID, err := s.jwt.ParseToken(token)
	if err != nil {
		return 0, uuid.Nil, err
	}

	err = s.redis.Delete(token)
	if err != nil {
		return 0, uuid.Nil, err
	}

	return userID, sessionID, nil
}

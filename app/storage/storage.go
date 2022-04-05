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

func (s *Storage) GetSession(token string) (*SessionData, error) {
	userID, sessionID, err := s.jwt.ParseToken(token)
	if err != nil {
		return nil, err
	}

	data, err := s.redis.Get(token)
	if err != nil {
		return nil, err
	}

	return &SessionData{
		ID:     sessionID,
		UserID: userID,
		Data:   data,
	}, nil
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

func (s *Storage) UpdateSession(token string, userData []byte) error {
	return s.redis.Set(token, userData)
}

func (s *Storage) DeleteSession(token string) error {
	return s.redis.Delete(token)
}

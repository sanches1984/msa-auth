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

func (s *Storage) DecodeToken(token string) (int64, uuid.UUID, error) {
	return s.jwt.ParseToken(token)
}

func (s *Storage) GetSessionData(token string) ([]byte, error) {
	return s.redis.Get(token)
}

func (s *Storage) GetSessionDataByUUID(sessionID uuid.UUID) ([]byte, error) {
	token, err := s.redis.Get(sessionID.String())
	if err != nil {
		return nil, err
	}
	return s.GetSessionData(string(token))
}

func (s *Storage) CreateSession(userID int64, userData []byte) (*Session, error) {
	sessionID := uuid.NewV4()
	session, err := s.createNewSession(userID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.redis.Set(session.Access.Value, userData); err != nil {
		return nil, err
	}

	if err := s.redis.Set(sessionID.String(), []byte(session.Access.Value)); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Storage) RefreshSession(userID int64, sessionID uuid.UUID, userData []byte) (*Session, error) {
	session, err := s.createNewSession(userID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.DeleteSessionByUUID(sessionID); err != nil {
		return nil, err
	}

	if err := s.redis.Set(session.Access.Value, userData); err != nil {
		return nil, err
	}

	if err := s.redis.Set(sessionID.String(), []byte(session.Access.Value)); err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Storage) UpdateSessionData(token string, userData []byte) error {
	return s.redis.Set(token, userData)
}

func (s *Storage) DeleteSession(token string) error {
	_, sessionID, err := s.jwt.ParseToken(token)
	if err != nil {
		return err
	}

	return s.deleteSessionRecords(sessionID, token)
}

func (s *Storage) DeleteSessionByUUID(sessionID uuid.UUID) error {
	token, err := s.redis.Get(sessionID.String())
	if err != nil {
		return err
	}

	return s.deleteSessionRecords(sessionID, string(token))
}

func (s *Storage) createNewSession(userID int64, sessionID uuid.UUID) (*Session, error) {
	access, err := s.jwt.NewAccessToken(userID, sessionID)
	if err != nil {
		return nil, err
	}
	refresh, err := s.jwt.NewRefreshToken(userID, sessionID)
	if err != nil {
		return nil, err
	}

	return &Session{
		ID:     sessionID,
		UserID: userID,
		Access: Token{
			Value:     access.Value,
			ExpiresIn: access.ExpiresAt,
		},
		Refresh: Token{
			Value:     refresh.Value,
			ExpiresIn: refresh.ExpiresAt,
		},
	}, nil
}

func (s *Storage) deleteSessionRecords(sessionID uuid.UUID, token string) error {
	// danger!
	if err := s.redis.Delete(token); err != nil {
		return err
	}
	return s.redis.Delete(sessionID.String())
}

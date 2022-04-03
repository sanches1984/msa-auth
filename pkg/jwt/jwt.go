package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type userJwt struct {
	UserID    int64 `json:"u"`
	SessionID int64 `json:"s"`
	jwt.StandardClaims
}

type Service struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewService(accessTTL, refreshTTL time.Duration, secret string) *Service {
	return &Service{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) NewAccessToken(userID, sessionID int64) (string, error) {
	return s.newToken(userID, sessionID, s.accessTTL)
}

func (s *Service) NewRefreshToken(userID, sessionID int64) (string, error) {
	return s.newToken(userID, sessionID, s.refreshTTL)
}

func (s *Service) ParseToken(token string) (int64, int64, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &userJwt{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("can't decode jwt token")
		}

		return s.secret, nil
	})
	if err != nil {
		return 0, 0, err
	}

	if jwtToken == nil || jwtToken.Claims == nil {
		return 0, 0, errors.New("token or claims is null")
	}

	authJwt, ok := jwtToken.Claims.(*userJwt)
	if !ok || !jwtToken.Valid {
		return 0, 0, errors.New("not valid token")
	}

	return authJwt.UserID, authJwt.SessionID, nil
}

func (s *Service) newToken(userID, sessionID int64, ttl time.Duration) (string, error) {
	now := time.Now()
	object := &userJwt{
		SessionID: sessionID,
		UserID:    userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.UTC().Unix(),
			ExpiresAt: now.Add(ttl).UTC().Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, object)
	return jwtToken.SignedString(s.secret)
}

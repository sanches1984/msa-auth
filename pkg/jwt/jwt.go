package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")

type userJwt struct {
	UserID    int64     `json:"u"`
	SessionID uuid.UUID `json:"s"`
	jwt.StandardClaims
}

type Token struct {
	Value     string
	ExpiresAt int32
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

func (s *Service) NewAccessToken(userID int64, sessionID uuid.UUID) (Token, error) {
	return s.newToken(userID, sessionID, s.accessTTL)
}

func (s *Service) NewRefreshToken(userID int64, sessionID uuid.UUID) (Token, error) {
	return s.newToken(userID, sessionID, s.refreshTTL)
}

func (s *Service) ParseToken(token string) (int64, uuid.UUID, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &userJwt{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("can't decode jwt token")
		}

		return s.secret, nil
	})
	if err != nil {
		return 0, uuid.Nil, err
	}

	if jwtToken == nil || jwtToken.Claims == nil {
		return 0, uuid.Nil, errors.New("token or claims is null")
	}

	authJwt, ok := jwtToken.Claims.(*userJwt)
	if !ok || !jwtToken.Valid {
		return 0, uuid.Nil, ErrInvalidToken
	}

	return authJwt.UserID, authJwt.SessionID, nil
}

func (s *Service) newToken(userID int64, sessionID uuid.UUID, ttl time.Duration) (Token, error) {
	now := time.Now()
	object := &userJwt{
		UserID:    userID,
		SessionID: sessionID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.UTC().Unix(),
			ExpiresAt: now.Add(ttl).UTC().Unix(),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, object)
	token, err := jwtToken.SignedString(s.secret)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Value:     token,
		ExpiresAt: int32(object.StandardClaims.ExpiresAt),
	}, nil
}

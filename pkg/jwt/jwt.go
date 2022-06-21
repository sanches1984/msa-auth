package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"
)

var ErrInvalidToken = errors.New("invalid token")
var ErrEmptyToken = errors.New("token or claims is null")

type userJwt struct {
	UserID    int64     `json:"u"`
	SessionID uuid.UUID `json:"s"`
	jwt.RegisteredClaims
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
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("can't decode jwt token")
		}

		return s.secret, nil
	})
	if err != nil {
		return 0, uuid.Nil, err
	}

	if jwtToken == nil || jwtToken.Claims == nil {
		return 0, uuid.Nil, ErrEmptyToken
	} else if !jwtToken.Valid {
		return 0, uuid.Nil, ErrInvalidToken
	} else if err := jwtToken.Claims.Valid(); err != nil {
		return 0, uuid.Nil, ErrInvalidToken
	}

	mc, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, uuid.Nil, ErrEmptyToken
	}

	userID, err := strconv.ParseInt(mc["jti"].(string), 10, 64)
	if err != nil {
		return 0, uuid.Nil, ErrEmptyToken
	}

	sessionID, err := uuid.FromString(mc["sub"].(string))
	if err != nil {
		return 0, uuid.Nil, ErrEmptyToken
	}

	return userID, sessionID, nil
}

func (s *Service) newToken(userID int64, sessionID uuid.UUID, ttl time.Duration) (Token, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		ID:        strconv.FormatInt(userID, 10),
		Subject:   sessionID.String(),
		IssuedAt:  &jwt.NumericDate{Time: now},
		ExpiresAt: &jwt.NumericDate{Time: now.Add(ttl)},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString(s.secret)
	if err != nil {
		return Token{}, err
	}

	return Token{
		Value:     token,
		ExpiresAt: int32(claims.ExpiresAt.Unix()),
	}, nil
}

package storage

import (
	"github.com/golang/mock/gomock"
	"github.com/sanches1984/auth/internal/pkg/storage/mocks"
	"github.com/sanches1984/auth/pkg/jwt"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StorageSuite struct {
	suite.Suite

	ctrl  *gomock.Controller
	redis *mocks.MockRedis
	jwt   *mocks.MockJwtService
}

func (s *StorageSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.redis = mocks.NewMockRedis(s.ctrl)
	s.jwt = mocks.NewMockJwtService(s.ctrl)
}

func (s *StorageSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestStorage(t *testing.T) {
	suite.Run(t, new(StorageSuite))
}

func (s *StorageSuite) TestDecodeToken() {
	token := "token"
	user := int64(123)
	session := uuid.NewV4()
	s.jwt.EXPECT().ParseToken(token).Return(user, session, nil).Times(1)

	userID, sessionID, err := New(s.redis, s.jwt).DecodeToken(token)
	s.NoError(err)
	s.Equal(user, userID)
	s.Equal(session, sessionID)
}

func (s *StorageSuite) TestGetSessionData() {
	token := "token"
	s.redis.EXPECT().Get(token).Return([]byte("hello"), nil).Times(1)

	data, err := New(s.redis, s.jwt).GetSessionData(token)
	s.NoError(err)
	s.Equal([]byte("hello"), data)
}

func (s *StorageSuite) TestGetSessionDataByUUID() {
	sessionID := uuid.NewV4()
	token := "token"
	s.redis.EXPECT().Get(sessionID.String()).Return([]byte(token), nil).Times(1)
	s.redis.EXPECT().Get(token).Return([]byte("hello"), nil).Times(1)

	data, err := New(s.redis, s.jwt).GetSessionDataByUUID(sessionID)
	s.NoError(err)
	s.Equal([]byte("hello"), data)
}

func (s *StorageSuite) TestCreateSession() {
	userID := int64(123)
	userData := []byte("hello")
	s.jwt.EXPECT().NewAccessToken(userID, gomock.Any()).Return(jwt.Token{Value: "token1", ExpiresAt: 111}, nil).Times(1)
	s.jwt.EXPECT().NewRefreshToken(userID, gomock.Any()).Return(jwt.Token{Value: "token2", ExpiresAt: 222}, nil).Times(1)
	s.redis.EXPECT().Set("token1", userData).Return(nil).Times(1)
	s.redis.EXPECT().Set(gomock.Any(), []byte("token1")).Return(nil).Times(1)

	session, err := New(s.redis, s.jwt).CreateSession(userID, userData)
	s.NoError(err)
	s.Equal(userID, session.UserID)
	s.Equal("token1", session.Access.Value)
	s.Equal("token2", session.Refresh.Value)
}

func (s *StorageSuite) TestRefreshSession() {
	userID := int64(123)
	sessionID := uuid.NewV4()
	userData := []byte("hello")
	s.jwt.EXPECT().NewAccessToken(userID, sessionID).Return(jwt.Token{Value: "token1", ExpiresAt: 111}, nil).Times(1)
	s.jwt.EXPECT().NewRefreshToken(userID, sessionID).Return(jwt.Token{Value: "token2", ExpiresAt: 222}, nil).Times(1)

	s.redis.EXPECT().Get(sessionID.String()).Return([]byte("old_token"), nil).Times(1)
	s.redis.EXPECT().Delete("old_token").Return(nil).Times(1)
	s.redis.EXPECT().Delete(sessionID.String()).Return(nil).Times(1)

	s.redis.EXPECT().Set("token1", userData).Return(nil).Times(1)
	s.redis.EXPECT().Set(sessionID.String(), []byte("token1")).Return(nil).Times(1)

	session, err := New(s.redis, s.jwt).RefreshSession(userID, sessionID, userData)
	s.NoError(err)
	s.Equal(userID, session.UserID)
	s.Equal(sessionID, session.ID)
	s.Equal("token1", session.Access.Value)
	s.Equal("token2", session.Refresh.Value)
}

func (s *StorageSuite) TestUpdateSessionData() {
	s.redis.EXPECT().Set("token", []byte("hello")).Return(nil).Times(1)
	err := New(s.redis, s.jwt).UpdateSessionData("token", []byte("hello"))
	s.NoError(err)
}

func (s *StorageSuite) TestDeleteSession() {
	token := "token"
	sessionID := uuid.NewV4()
	s.jwt.EXPECT().ParseToken(token).Return(int64(123), sessionID, nil).Times(1)
	s.redis.EXPECT().Delete(token).Return(nil).Times(1)
	s.redis.EXPECT().Delete(sessionID.String()).Return(nil).Times(1)

	err := New(s.redis, s.jwt).DeleteSession(token)
	s.NoError(err)
}

func (s *StorageSuite) TestDeleteSessionByUUID() {
	token := "token"
	sessionID := uuid.NewV4()
	s.redis.EXPECT().Get(sessionID.String()).Return([]byte(token), nil).Times(1)
	s.redis.EXPECT().Delete(token).Return(nil).Times(1)
	s.redis.EXPECT().Delete(sessionID.String()).Return(nil).Times(1)

	err := New(s.redis, s.jwt).DeleteSessionByUUID(sessionID)
	s.NoError(err)
}

package service

import (
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/internal/app/service/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	repo    *mocks.MockRepository
	storage *mocks.MockStorage
	logger  zerolog.Logger
}

func (s *AuthSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockRepository(s.ctrl)
	s.storage = mocks.NewMockStorage(s.ctrl)
	s.logger = zerolog.Nop()
}

func (s *AuthSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}

func (s *AuthSuite) TestLogin_Success() {
	// todo
}

func (s *AuthSuite) TestLogin_Error() {
	// todo
}

func (s *AuthSuite) TestLogout_Success() {
	// todo
}

func (s *AuthSuite) TestLogout_Error() {
	// todo
}

func (s *AuthSuite) TestChangePassword_Success() {
	// todo
}

func (s *AuthSuite) TestChangePassword_Error() {
	// todo
}

func (s *AuthSuite) TestNewAccessTokenByRefreshToken_Success() {
	// todo
}

func (s *AuthSuite) TestNewAccessTokenByRefreshToken_Error() {
	// todo
}

func (s *AuthSuite) TestValidateToken_Success() {
	// todo
}

func (s *AuthSuite) TestValidateToken_Error() {
	// todo
}

func (s *AuthSuite) TestUpdateSessionData_Success() {
	// todo
}

func (s *AuthSuite) TestUpdateSessionData_Error() {
	// todo
}

func (s *AuthSuite) TestGetUserSessions_Success() {
	// todo
}

func (s *AuthSuite) TestGetUserSessions_Error() {
	// todo
}

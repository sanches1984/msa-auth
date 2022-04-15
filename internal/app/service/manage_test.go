package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/sanches1984/auth/internal/app/model"
	"github.com/sanches1984/auth/internal/app/service/mocks"
	errs "github.com/sanches1984/auth/pkg/errors"
	api "github.com/sanches1984/auth/proto/api"
	"github.com/sanches1984/gopkg-pg-orm/pager"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ManageSuite struct {
	suite.Suite

	ctrl    *gomock.Controller
	repo    *mocks.MockRepository
	storage *mocks.MockStorage
	logger  zerolog.Logger
}

func (s *ManageSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mocks.NewMockRepository(s.ctrl)
	s.storage = mocks.NewMockStorage(s.ctrl)
	s.logger = zerolog.Nop()
}

func (s *ManageSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestManageService(t *testing.T) {
	suite.Run(t, new(ManageSuite))
}

func (s *ManageSuite) TestCreateUser_Success() {
	ctx := context.Background()

	s.repo.EXPECT().CreateUser(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *model.User) error {
		user.ID = 123
		return nil
	}).Times(1)

	resp, err := NewManageService(s.repo, s.storage, s.logger).CreateUser(ctx, &api.CreateUserRequest{
		Login:    "login",
		Password: "password",
	})

	s.NoError(err)
	s.Equal(&api.CreateUserResponse{UserId: 123}, resp)
}

func (s *ManageSuite) TestCreateUser_Error() {
	ctx := context.Background()
	repoErr := errors.New("some error")
	s.repo.EXPECT().CreateUser(ctx, gomock.Any()).Return(repoErr).Times(1)

	resp, err := NewManageService(s.repo, s.storage, s.logger).CreateUser(ctx, &api.CreateUserRequest{
		Login:    "login",
		Password: "password",
	})

	s.Nil(resp)
	s.EqualError(err, repoErr.Error())
}

func (s *ManageSuite) TestDeleteUser_Success() {
	ctx := context.Background()
	sessionID1 := uuid.NewV4()
	sessionID2 := uuid.NewV4()
	user := &model.User{ID: 123, Login: "login", PasswordHash: "hash"}
	tokens := model.RefreshTokenList{&model.RefreshToken{SessionID: sessionID1}, &model.RefreshToken{SessionID: sessionID2}}

	s.repo.EXPECT().GetUser(ctx, model.UserFilter{ID: 123}).Return(user, nil).Times(1)
	s.repo.EXPECT().GetRefreshTokens(ctx, model.RefreshTokenFilter{UserID: user.ID}).Return(tokens, nil).Times(1)
	s.storage.EXPECT().DeleteSessionByUUID(sessionID1).Return(nil).Times(1)
	s.storage.EXPECT().DeleteSessionByUUID(sessionID2).Return(nil).Times(1)
	s.repo.EXPECT().DeleteRefreshToken(ctx, model.RefreshTokenFilter{UserID: user.ID}).Times(1)
	s.repo.EXPECT().DeleteUser(ctx, user).Return(nil).Times(1)

	resp, err := NewManageService(s.repo, s.storage, s.logger).DeleteUser(ctx, &api.DeleteUserRequest{UserId: 123})
	s.NoError(err)
	s.Equal(&api.DeleteUserResponse{SessionId: []string{sessionID1.String(), sessionID2.String()}}, resp)
}

func (s *ManageSuite) TestDeleteUser_Error() {
	ctx := context.Background()

	s.repo.EXPECT().GetUser(ctx, model.UserFilter{ID: 123}).Return(nil, nil).Times(1)

	resp, err := NewManageService(s.repo, s.storage, s.logger).DeleteUser(ctx, &api.DeleteUserRequest{UserId: 123})
	s.Nil(resp)
	s.EqualError(err, errs.ErrUserNotFound.Error())
}

func (s *ManageSuite) TestGetUsers_Success() {
	ctx := context.Background()
	now := time.Now()
	users := model.UserList{
		&model.User{ID: 1, Login: "login1", Created: now, Updated: now},
		&model.User{ID: 2, Login: "login2", Created: now, Updated: now},
	}

	s.repo.EXPECT().GetUsers(ctx, model.UserFilter{Order: model.UserOrderLoginDesc}, pager.NewPagerWithPageSize(2, 5)).Return(users, nil)

	resp, err := NewManageService(s.repo, s.storage, s.logger).GetUsers(ctx, &api.GetUsersRequest{
		Order:    api.GetUsersRequest_LOGIN_DESC,
		Page:     2,
		PageSize: 5,
	})
	s.NoError(err)
	s.Equal(&api.GetUsersResponse{Users: []*api.User{
		{
			Id:      1,
			Login:   "login1",
			Created: now.Format(time.RFC3339),
			Updated: now.Format(time.RFC3339),
		},
		{
			Id:      2,
			Login:   "login2",
			Created: now.Format(time.RFC3339),
			Updated: now.Format(time.RFC3339),
		},
	}}, resp)
}

func (s *ManageSuite) TestGetUsers_Error() {
	ctx := context.Background()
	dbErr := errors.New("internal")

	s.repo.EXPECT().GetUsers(ctx, model.UserFilter{}, pager.NewPagerWithPageSize(0, 0)).Return(nil, dbErr)

	resp, err := NewManageService(s.repo, s.storage, s.logger).GetUsers(ctx, &api.GetUsersRequest{})
	s.Nil(resp)
	s.EqualError(err, dbErr.Error())
}

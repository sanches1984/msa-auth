package service

import (
	"context"
	"github.com/rs/zerolog"
	api "github.com/sanches1984/auth/proto/api"
)

type ManageService struct {
	api.ManageServiceServer

	factory Factory
	logger  zerolog.Logger
}

func NewManageService(factory Factory, logger zerolog.Logger) *ManageService {
	return &ManageService{
		factory: factory,
		logger:  logger,
	}
}

func (s *ManageService) CreateUser(ctx context.Context, r *api.CreateUserRequest) (*api.CreateUserResponse, error) {
	// todo
	return &api.CreateUserResponse{Id: 12}, nil
}

func (s *ManageService) DeleteUser(ctx context.Context, r *api.DeleteUserRequest) (*api.DeleteUserResponse, error) {
	// todo
	return &api.DeleteUserResponse{Deleted: true}, nil
}

func (s *ManageService) GetUserList(ctx context.Context, r *api.GetUserListRequest) (*api.GetUserListResponse, error) {
	// todo
	return &api.GetUserListResponse{}, nil
}

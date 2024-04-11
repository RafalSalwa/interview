package services

import (
	"context"
	"strconv"

	"github.com/RafalSalwa/interview-app-srv/cmd/auth_service/internal/repository"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
)

type UserService interface {
	Load(ctx context.Context, id string) (*models.UserDBModel, error)
}

func NewORMUserService(repository repository.UserRepository) UserService {
	return &userService{repository: repository}
}

type userService struct {
	repository repository.UserRepository
}

func (s *userService) Load(ctx context.Context, id string) (*models.UserDBModel, error) {
	uid, _ := strconv.ParseInt(id, 10, 64)
	res, err := s.repository.GetOrCreate(ctx, uid)
	if err != nil {
		return nil, err
	}
	return res, err
}

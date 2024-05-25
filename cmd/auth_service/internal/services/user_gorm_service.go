package services

import (
	"context"
	"strconv"

	"github.com/RafalSalwa/auth-api/cmd/auth_service/internal/repository"
	"github.com/RafalSalwa/auth-api/pkg/models"
)

type (
	UserService interface {
		Load(ctx context.Context, id string) (*models.UserDBModel, error)
	}
	userService struct {
		repo repository.UserRepository
	}
)

func NewORMUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Load(ctx context.Context, id string) (*models.UserDBModel, error) {
	uid, _ := strconv.ParseInt(id, 10, 64)
	res, err := s.repo.GetOrCreate(ctx, uid)
	if err != nil {
		return nil, err
	}
	return res, err
}

package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
)

type UserRepository interface {
	FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error)
	Save(ctx context.Context, user *models.UserDBModel) error
	Update(ctx context.Context, user *models.UserDBModel) error

	SingUp(user *models.UserDBModel) error
	Load(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error)
	ById(ctx context.Context, id int64) (*models.UserDBModel, error)
	ByLogin(ctx context.Context, user *models.SignInUserRequest) (*models.UserDBModel, error)
	ConfirmVerify(ctx context.Context, vCode string) error
	UpdateLastLogin(ctx context.Context, u *models.UserDBModel) (*models.UserDBModel, error)
	FindUserByID(uid int64) (*models.UserDBModel, error)
	ChangePassword(ctx context.Context, userid int64, password string) error
	BeginTx() *gorm.DB
	GetConnection() *gorm.DB
}

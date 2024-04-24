package repository

import (
	"context"
	"errors"
	"github.com/RafalSalwa/interview-app-srv/pkg/hashing"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"go.opentelemetry.io/otel"
	"gorm.io/gorm"
	"time"
)

type UserAdapter struct {
	DB *gorm.DB
}

func (r *UserAdapter) Save(ctx context.Context, user *models.UserDBModel) error {
	// TODO implement me
	panic("implement me")
}

func (r *UserAdapter) Update(ctx context.Context, user *models.UserDBModel) error {
	_, span := otel.GetTracerProvider().Tracer("repository").Start(ctx, "Repository/Update")
	defer span.End()

	res := r.DB.Updates(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *UserAdapter) FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	_, span := otel.GetTracerProvider().Tracer("repository").Start(ctx, "Repository/FindOne")
	defer span.End()

	if err := r.DB.Where(&user).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}

	if user.Email == "" || user.Id == 0 {
		return nil, nil
	}
	return user, nil
}

func NewUserAdapter(db *gorm.DB) UserRepository {
	return &UserAdapter{DB: db}
}

func (r *UserAdapter) ChangePassword(ctx context.Context, userid int64, password string) error {
	user := models.UserDBModel{Id: userid}
	return r.DB.Model(user).
		Updates(models.UserDBModel{
			Password: password,
			Active:   true,
		}).
		Error
}

func (r *UserAdapter) Load(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	_, span := otel.GetTracerProvider().Tracer("repository").Start(ctx, "Repository/Load")
	defer span.End()

	if err := r.DB.Where(&user).Limit(1).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserAdapter) ConfirmVerify(ctx context.Context, vCode string) error {
	user := models.UserDBModel{VerificationCode: vCode}
	_, err := r.Load(ctx, &user)
	if err != nil {
		return err
	}
	var count int64

	r.DB.Where(&user).First(&user).Count(&count)
	if count == 0 {
		return errors.New("NotFound")
	}
	if count == 1 && user.Active && user.Verified {
		return errors.New("AlreadyActivated")
	}
	return r.DB.Model(user).Where(&user).
		Updates(models.UserDBModel{
			Verified: true,
			Active:   true,
		}).
		Error
}

func (r *UserAdapter) SingUp(user *models.UserDBModel) error {
	return r.DB.Create(&user).Error
}

func (r *UserAdapter) ById(_ context.Context, id int64) (*models.UserDBModel, error) {
	var user models.UserDBModel
	r.DB.First(&user, "id = ?", id)
	return &user, nil
}

func (r *UserAdapter) ByLogin(_ context.Context, user *models.SignInUserRequest) (*models.UserDBModel, error) {
	var dbUser models.UserDBModel

	r.DB.First(&dbUser, "username = ? OR email = ?", user.Username, user.Email)
	if hashing.CheckPasswordHash(user.Password, dbUser.Password) {
		return &dbUser, nil
	}
	return nil, nil
}

func (r *UserAdapter) UpdateLastLogin(ctx context.Context, u *models.UserDBModel) (*models.UserDBModel, error) {
	now := time.Now()
	r.DB.Model(u).Update("LastLogin", now)
	u.LastLogin = &now
	return u, nil
}
func (r *UserAdapter) BeginTx() *gorm.DB {
	return r.DB.Begin().Begin()
}
func (r *UserAdapter) GetConnection() *gorm.DB {
	return r.DB
}

func (r *UserAdapter) FindUserByID(uid int64) (*models.UserDBModel, error) {
	return nil, nil
}

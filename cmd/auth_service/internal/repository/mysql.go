package repository

import (
	"context"
	"errors"
	"time"

	"github.com/RafalSalwa/interview-app-srv/pkg/hashing"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"gorm.io/gorm"
)

type UserAdapter struct {
	DB *gorm.DB
}

func (r *UserAdapter) Exists(ctx context.Context, udb *models.UserDBModel) bool {
	_, span := otel.GetTracerProvider().Tracer("MySQL").Start(ctx, "Repository/Exists")
	defer span.End()

	err := r.DB.Where(udb).First(udb).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
	}
	if udb.Id != 0 {
		return true
	}

	return false
}

func newMySQLUserRepository(db *gorm.DB) UserRepository {
	return &UserAdapter{DB: db}
}

func (r *UserAdapter) FindOne(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	_, span := otel.GetTracerProvider().Tracer("auth-handler").Start(ctx, "Handler SignUpUser")
	defer span.End()

	if err := r.DB.Where(&user).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserAdapter) Save(ctx context.Context, user *models.UserDBModel) error {
	_, span := otel.GetTracerProvider().Tracer("auth-handler").Start(ctx, "Handler SignUpUser")
	defer span.End()

	res := r.DB.Save(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *UserAdapter) Update(ctx context.Context, user models.UserDBModel) error {
	_, span := otel.GetTracerProvider().Tracer("auth-handler").Start(ctx, "Handler SignUpUser")
	defer span.End()

	res := r.DB.Updates(&user)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (r *UserAdapter) FindAll(ctx context.Context, user *models.UserDBModel) ([]models.UserDBModel, error) {
	_, span := otel.GetTracerProvider().Tracer("auth-handler").Start(ctx, "Handler SignUpUser")
	defer span.End()

	var all []models.UserDBModel
	res := r.DB.Find(&all)
	if res.Error != nil {
		return nil, res.Error
	}

	return all, nil
}

func (r *UserAdapter) Confirm(ctx context.Context, user *models.UserDBModel) error {
	_, span := otel.GetTracerProvider().Tracer("auth_service-repository").Start(ctx, "Confirm")
	defer span.End()

	return r.DB.Model(user).Updates(models.UserDBModel{
		Verified: true,
		Active:   true,
	}).Error
}

func (r *UserAdapter) SignUp(ctx context.Context, user models.UserDBModel) error {
	_, span := otel.GetTracerProvider().Tracer("auth_service-repository").Start(ctx, "MySQL Repository SignUp")
	defer span.End()

	if err := r.DB.Create(&user).Error; err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		return err
	}
	return nil
}

func (r *UserAdapter) GetOrCreate(ctx context.Context, id int64) (*models.UserDBModel, error) {
	_, span := otel.GetTracerProvider().Tracer("auth_service-repository").Start(ctx, "GetOrCreate")
	defer span.End()
	var user models.UserDBModel
	r.DB.First(&user, "id = ?", id)
	return &user, nil
}

func (r *UserAdapter) ByLogin(ctx context.Context, user *models.SignInUserRequest) (*models.UserDBModel, error) {
	var dbUser models.UserDBModel

	r.DB.First(&dbUser, "username = ? OR email = ?", user.Username, user.Email)
	if hashing.CheckPasswordHash(user.Password, dbUser.Password) {
		return &dbUser, nil
	}
	return nil, nil
}

func (r *UserAdapter) UpdateLastLogin(ctx context.Context, u *models.UserDBModel) error {
	_, span := otel.GetTracerProvider().Tracer("auth_service-repository").Start(ctx, "UpdateLastLogin")
	defer span.End()

	now := time.Now()
	r.DB.Model(u).Update("LastLogin", now)
	u.LastLogin = &now
	return nil
}

func (r *UserAdapter) GetConnection() *gorm.DB {
	return r.DB
}

func (r *UserAdapter) FindUserById(uid int64) (*models.UserDBModel, error) {
	return nil, nil
}

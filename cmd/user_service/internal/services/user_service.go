package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/RafalSalwa/auth-api/pkg/encdec"
	"github.com/RafalSalwa/auth-api/pkg/hashing"
	"github.com/RafalSalwa/auth-api/pkg/tracing"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RafalSalwa/auth-api/cmd/user_service/config"
	"github.com/RafalSalwa/auth-api/cmd/user_service/internal/repository"
	"github.com/RafalSalwa/auth-api/pkg/jwt"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/models"
	"github.com/RafalSalwa/auth-api/pkg/rabbitmq"
)

type UserServiceImpl struct {
	repository      repository.UserRepository
	rabbitPublisher *rabbitmq.Publisher
	logger          *logger.Logger
	config          jwt.JWTConfig
}

type UserService interface {
	Find(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error)
	GetUser(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error)
	GetByID(ctx context.Context, id int64) (*models.UserDBModel, error)
	UsernameInUse(ctx context.Context, user *models.UserDBModel) (bool, error)
	StoreVerificationData(ctx context.Context, vCode string) error
	UpdateUser(user *models.UpdateUserRequest) (err error)
	LoginUser(user *models.SignInUserRequest) (*models.UserResponse, error)
	UpdateUserPassword(ctx context.Context, userid int64, password string) error
	CreateUser(user *models.SignUpUserRequest) (*models.UserResponse, error)
	GetByToken(ctx context.Context, token string) (*models.UserDBModel, error)
}

func NewUserService(ctx context.Context, cfg *config.Config, log *logger.Logger) UserServiceImpl {
	userRepository, errR := repository.NewUserRepository(ctx, cfg.App.RepositoryType, cfg)
	if errR != nil {
		log.Error().Err(errR).Msg("user:service:new")
	}

	publisher, errP := rabbitmq.NewPublisher(cfg.Rabbit)
	if errP != nil {
		log.Error().Err(errP).Msg("rabbitmq")
	}

	return UserServiceImpl{
		repository:      userRepository,
		rabbitPublisher: publisher,
		logger:          log,
		config:          cfg.JWTToken,
	}
}

func (s *UserServiceImpl) Find(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("service").Start(ctx, "Service/GetUser")
	defer span.End()

	udb, err := s.repository.FindOne(ctx, user)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	return udb, nil
}
func (s *UserServiceImpl) GetUser(ctx context.Context, user *models.UserDBModel) (*models.UserDBModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("service").Start(ctx, "Service/GetUser")
	defer span.End()

	find := &models.UserDBModel{Email: encdec.Encrypt(user.Email)}

	udb, err := s.repository.FindOne(ctx, find)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	_, err = hashing.Argon2IDComparePasswordAndHash(user.Password, udb.Password)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	return udb, nil
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id int64) (*models.UserDBModel, error) {
	user, err := s.repository.ById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) UsernameInUse(ctx context.Context, user *models.UserDBModel) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer("service").Start(ctx, "Service/UserExists")
	defer span.End()

	ur, err := s.repository.Load(ctx, user)
	if err != nil {
		return false, err
	}
	return ur == nil, nil
}

func (s *UserServiceImpl) StoreVerificationData(ctx context.Context, vCode string) error {
	ctx, span := otel.GetTracerProvider().Tracer("service").Start(ctx, "Service/StoreVerificationData")
	defer span.End()

	userDbModel := &models.UserDBModel{
		VerificationCode: vCode,
	}

	udb, err := s.repository.FindOne(ctx, userDbModel)
	if err != nil {
		tracing.RecordError(span, err)
		return err
	}
	if udb == nil {
		errUser := fmt.Errorf("user not found. user with verification code %s was not found", vCode)
		tracing.RecordError(span, errUser)
		return errUser
	}
	udb.Active = true
	udb.Verified = true
	err = s.repository.Update(ctx, udb)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) GetByToken(ctx context.Context, token string) (*models.UserDBModel, error) {
	ctx, span := otel.GetTracerProvider().Tracer("service").Start(ctx, "Service/StoreVerificationData")
	defer span.End()

	decodedToken, err := jwt.ValidateToken(token, s.config.Access.PublicKey)
	if err != nil {
		decodedToken, err = jwt.ValidateToken(token, s.config.Refresh.PublicKey)
	}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	uid, err := strconv.ParseInt(decodedToken, 10, 0)
	if err != nil {
		return nil, err
	}
	udb, err := s.GetByID(ctx, uid)
	if err != nil {
		return nil, err
	}
	udb.Email, _ = encdec.Decrypt(udb.Email)
	return udb, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.UpdateUserRequest) (err error) {
	// TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) LoginUser(user *models.SignInUserRequest) (*models.UserResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (s *UserServiceImpl) UpdateUserPassword(ctx context.Context, userid int64, password string) error {
	err := s.repository.ChangePassword(ctx, userid, password)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) CreateUser(user *models.SignUpUserRequest) (*models.UserResponse, error) {
	panic("implement me")
}

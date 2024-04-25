package services

import (
	"context"

	"github.com/RafalSalwa/auth-api/pkg/encdec"
	"github.com/RafalSalwa/auth-api/pkg/tracing"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RafalSalwa/auth-api/cmd/auth_service/config"
	"github.com/RafalSalwa/auth-api/cmd/auth_service/internal/repository"
	"github.com/RafalSalwa/auth-api/pkg/generator"
	"github.com/RafalSalwa/auth-api/pkg/hashing"
	"github.com/RafalSalwa/auth-api/pkg/jwt"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/models"
	"github.com/RafalSalwa/auth-api/pkg/rabbitmq"
	"go.opentelemetry.io/otel"
)

type AuthServiceImpl struct {
	repository      repository.UserRepository
	rabbitPublisher *rabbitmq.Publisher
	logger          *logger.Logger
	config          jwt.JWTConfig
}

func NewAuthService(ctx context.Context, cfg *config.Config, log *logger.Logger) AuthService {
	publisher, errP := rabbitmq.NewPublisher(cfg.Rabbit)
	if errP != nil {
		log.Error().Err(errP).Msg("auth:service:publisher")
		return nil
	}

	userRepository, errR := repository.NewUserRepository(ctx, cfg.App.RepositoryType, cfg)
	if errR != nil {
		log.Error().Err(errR).Msg("auth:service:repository")
		return nil
	}

	return &AuthServiceImpl{
		repository:      userRepository,
		rabbitPublisher: publisher,
		logger:          log,
		config:          cfg.JWTToken,
	}
}

func (a *AuthServiceImpl) SignUpUser(ctx context.Context, cur models.SignUpUserRequest) (*models.UserResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer("auth_service-service").Start(ctx, "Service/SignUpUser")
	defer span.End()

	udb := &models.UserDBModel{
		Email: encdec.Encrypt(cur.Email),
	}

	ok := a.repository.Exists(ctx, udb)
	if ok {
		return nil, status.Errorf(codes.AlreadyExists, "User with such credentials already exists")
	}
	if err := hashing.Validate(cur.Password, cur.PasswordConfirm); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	if err := udb.FromCreateUserReq(cur, true); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	udb.Password = hashing.Argon2ID(udb.Password)
	vcode, _ := generator.RandomString(64)
	udb.VerificationCode = vcode
	udb.Roles = map[string]interface{}{
		"Roles": struct {
			Role string
		}{"ROLE_USER"},
	}
	if errDB := a.repository.Save(ctx, udb); errDB != nil {
		return nil, status.Errorf(codes.Internal, errDB.Error())
	}
	if err := a.rabbitPublisher.Publish(ctx, udb.AMQP()); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	ur := &models.UserResponse{}
	err := ur.FromDBModel(udb)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	return ur, nil
}

func (a *AuthServiceImpl) SignInUser(ctx context.Context, reqUser *models.SignInUserRequest) (*models.UserResponse, error) {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "Service/SignInUser")
	defer span.End()

	udb := &models.UserDBModel{
		Email:    encdec.Encrypt(reqUser.Email),
		Verified: true,
		Active:   true,
	}
	udb, err := a.repository.FindOne(ctx, udb)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}
	if udb == nil {
		return nil, status.Errorf(codes.NotFound, "User with such credentials was not found")
	}
	if _, err = hashing.Argon2IDComparePasswordAndHash(reqUser.Password, udb.Password); err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	ur := &models.UserResponse{}
	err = ur.FromDBModel(udb)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	tp, err := jwt.GenerateTokenPair(a.config, udb.Id)
	if err != nil {
		tracing.RecordError(span, err)
		return nil, err
	}

	ur.AssignTokenPair(tp)
	return ur, nil
}

func (a *AuthServiceImpl) GetVerificationKey(ctx context.Context, email string) (*models.UserResponse, error) {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "GetVerificationKey")
	defer span.End()
	user := &models.UserDBModel{
		Email: encdec.Encrypt(email),
	}
	dbUser, err := a.repository.FindOne(ctx, user)
	if err != nil {
		return nil, err
	}
	ur := &models.UserResponse{}
	err = ur.FromDBModel(dbUser)
	if err != nil {
		return nil, err
	}
	return ur, nil
}

func (a *AuthServiceImpl) Find(ctx context.Context, user *models.UserDBModel) (*models.UserResponse, error) {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "FindAll")
	defer span.End()
	dbUser, err := a.repository.FindOne(ctx, user)
	if err != nil {
		return nil, err
	}
	if dbUser == nil {
		return nil, nil
	}

	ur := &models.UserResponse{}
	err = ur.FromDBModel(dbUser)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

func (a *AuthServiceImpl) Load(ctx context.Context, user *models.UserDBModel) (*models.UserResponse, error) {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "Service/FindOne")
	defer span.End()

	dbUser, err := a.repository.FindOne(ctx, user)
	if err != nil {
		return nil, err
	}
	if dbUser == nil {
		return nil, nil
	}
	err = a.repository.Update(ctx, *dbUser)
	if err != nil {
		return nil, err
	}

	tp, err := jwt.GenerateTokenPair(a.config, dbUser.Id)
	_, _ = jwt.DecodeToken(tp.AccessToken, a.config.Access.PublicKey)
	if err != nil {
		a.logger.Error().Err(err).Msg("token_pair")
		return nil, err
	}

	ur := &models.UserResponse{}
	err = ur.FromDBModel(dbUser)
	if err != nil {
		return nil, err
	}
	ur.AssignTokenPair(tp)

	return ur, nil
}

func (a *AuthServiceImpl) Verify(ctx context.Context, vCode string) error {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "Verify")
	defer span.End()

	udb := &models.UserDBModel{
		VerificationCode: vCode,
	}
	dbUser, err := a.repository.FindOne(ctx, udb)
	if err != nil {
		return err
	}

	if errV := a.repository.Confirm(ctx, udb); errV != nil {
		return errV
	}
	ur := &models.UserResponse{}

	return ur.FromDBModel(dbUser)
}

func (a *AuthServiceImpl) FindUserByID(uid int64) (*models.UserDBModel, error) {
	ctx := context.Background()
	dbUser, err := a.repository.GetOrCreate(ctx, uid)
	if err != nil {
		return nil, err
	}
	return dbUser, nil
}

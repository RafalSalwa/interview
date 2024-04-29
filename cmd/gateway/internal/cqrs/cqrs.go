package cqrs

import (
	"context"

	"github.com/RafalSalwa/auth-api/cmd/gateway/config"
	"github.com/RafalSalwa/auth-api/cmd/gateway/internal/cqrs/command"
	"github.com/RafalSalwa/auth-api/cmd/gateway/internal/cqrs/query"
	"github.com/RafalSalwa/auth-api/cmd/gateway/internal/rpc_client"
	"github.com/RafalSalwa/auth-api/pkg/models"
	intrvproto "github.com/RafalSalwa/auth-api/proto/grpc"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	SignUp           command.SignUpHandler
	SignIn           command.SignInHandler
	SignInByCode     command.SignInByCodeHandler
	ChangePassword   command.ChangePasswordHandler
	VerifyUserByCode command.VerifyCodeHandler
}

type Queries struct {
	GetUser     query.GetUserHandler
	UserDetails query.UserDetailsHandler

	VerificationCode query.VerificationCodeHandler
	FetchUser        query.FetchUserHandler
	UserExists       query.UserExistsHandler
	GetUserByCode    query.GetUserByCodeHandler
}

func NewService(cfg config.Grpc) (*Application, error) {
	authClient, err := rpc_client.NewAuthClient(cfg.AuthServicePort)
	if err != nil {
		return nil, err
	}

	userClient, err := rpc_client.NewUserClient(cfg.UserServicePort)
	if err != nil {
		return nil, err
	}

	return newApplication(authClient, userClient), nil
}

func newApplication(authClient intrvproto.AuthServiceClient, userClient intrvproto.UserServiceClient) *Application {
	return &Application{
		Commands: Commands{
			SignUp:           command.NewSignUpHandler(authClient),
			SignIn:           command.NewSignInHandler(authClient),
			SignInByCode:     command.NewSignInByCodeHandler(authClient),
			ChangePassword:   command.NewChangePasswordHandler(userClient),
			VerifyUserByCode: command.NewVerifyCodeHandler(userClient),
		},
		Queries: Queries{
			UserExists:       query.NewUserExistsHandler(userClient),
			UserDetails:      query.NewUserDetailsHandler(userClient),
			GetUser:          query.NewGetUserHandler(userClient),
			VerificationCode: query.NewVerificationCodeHandler(authClient),
			FetchUser:        query.NewFetchUserHandler(userClient),
			GetUserByCode:    query.NewGetUserByCodeHandler(userClient),
		},
	}
}

func (app *Application) CheckUserExistsQuery(ctx context.Context, email string) (bool, error) {
	exists, err := app.Queries.UserExists.Handle(ctx, email)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (app *Application) SignupUserCommand(ctx context.Context, req models.SignUpUserRequest) error {
	return app.Commands.SignUp.Handle(ctx, req)
}

func (app *Application) SigninCommand(ctx context.Context, req models.SignInUserRequest) (*models.UserResponse, error) {
	return app.Commands.SignIn.Handle(ctx, req)
}

func (app *Application) SigninByCodeCommand(ctx context.Context, email, authCode string) (*models.UserResponse, error) {
	return app.Commands.SignInByCode.Handle(ctx, email, authCode)
}

func (app *Application) GetVerificationCode(ctx context.Context, email string) (models.UserResponse, error) {
	return app.Queries.VerificationCode.Handle(ctx, email)
}

func (app *Application) FetchUser(ctx context.Context, email, password string) (models.UserResponse, error) {
	u := models.SignInUserRequest{
		Email:    email,
		Password: password,
	}
	return app.Queries.FetchUser.Handle(ctx, query.FetchUser{SignInUserRequest: u})
}

func (app *Application) ChangePassword(ctx context.Context, req *models.ChangePasswordRequest) error {
	return app.Commands.ChangePassword.Handle(ctx, command.ChangePassword{
		Id:       req.Id,
		Password: req.Password,
	})
}

func (app *Application) GetUser(ctx context.Context, id models.UserRequest) (models.UserResponse, error) {
	ur, err := app.Queries.GetUser.Handle(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	user, err := app.SigninByCodeCommand(ctx, ur.Email, ur.VerificationCode)
	if err != nil {
		return models.UserResponse{}, err
	}
	ur.AccessToken = user.AccessToken
	ur.RefreshToken = user.RefreshToken

	return ur, nil
}

func (app *Application) GetUserByCode(ctx context.Context, vCode string) (models.UserResponse, error) {
	return app.Queries.GetUserByCode.Handle(ctx, vCode)
}

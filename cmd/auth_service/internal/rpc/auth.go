package rpc

import (
	"context"
	"errors"

	"github.com/RafalSalwa/interview-app-srv/pkg/encdec"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"github.com/RafalSalwa/interview-app-srv/pkg/tracing"
	pb "github.com/RafalSalwa/interview-app-srv/proto/grpc"
	"go.opentelemetry.io/otel"
	otelcodes "go.opentelemetry.io/otel/codes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

func (a *Auth) SignInUser(ctx context.Context, req *pb.SignInUserInput) (*pb.SignInUserResponse, error) {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "GRPC SignInUser")
	defer span.End()

	loginUser := &models.SignInUserRequest{
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	ur, err := a.authService.SignInUser(ctx, loginUser)
	if err != nil {
		a.logger.Error().Err(err).Msg("rpc:service:signin")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.SignInUserResponse{
		AccessToken:  ur.AccessToken,
		RefreshToken: ur.RefreshToken,
	}
	return res, nil
}
func (a *Auth) SignInByCode(ctx context.Context, req *pb.SignInByCodeUserInput) (*pb.SignInUserResponse, error) {
	ctx, span := tracing.InitSpan(ctx, "auth_service-rpc", "GRPC SignInByCode")
	defer span.End()

	loginUser := &models.UserDBModel{
		Email:            encdec.Encrypt(req.GetEmail()),
		VerificationCode: req.GetAuthCode(),
	}
	ur, err := a.authService.Load(ctx, loginUser)
	if err != nil {
		a.logger.Error().Err(err).Msg("rpc:service:signin")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.SignInUserResponse{
		AccessToken:  ur.AccessToken,
		RefreshToken: ur.RefreshToken,
	}
	return res, nil
}

func (a *Auth) SignUpUser(ctx context.Context, req *pb.SignUpUserInput) (*pb.SignUpUserResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer("grpc func").Start(ctx, "RPC/SignUpUser")
	defer span.End()

	userSignUp := models.SignUpUserRequest{
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
	ur, err := a.authService.SignUpUser(ctx, userSignUp)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(otelcodes.Error, err.Error())
		a.logger.Error().Err(err).Msg("rpc:signup")
		if s, ok := status.FromError(err); ok {
			return nil, status.Errorf(s.Code(), s.Message())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	res := &pb.SignUpUserResponse{
		Id:                ur.Id,
		Username:          ur.Username,
		VerificationToken: ur.VerificationCode,
		CreatedAt:         timestamppb.New(ur.CreatedAt),
	}
	return res, nil
}

func (a *Auth) GetVerificationKey(
	ctx context.Context,
	in *pb.VerificationCodeRequest) (*pb.VerificationCodeResponse, error) {
	ur, err := a.authService.GetVerificationKey(ctx, in.GetEmail())
	if err != nil {
		a.logger.Error().Err(err).Msg("rpc:service:getkey")
		if err.Error() == "record not found" {
			return nil, status.Errorf(codes.NotFound, "user with such credentials was not found")
		}
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return &pb.VerificationCodeResponse{Code: ur.VerificationCode}, nil
}

package command

import (
	"context"

	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type VerifyCode struct {
	VerificationCode string
}

type VerifyCodeHandler struct {
	grpcUser intrvproto.UserServiceClient
}

func NewVerifyCodeHandler(grpcUser intrvproto.UserServiceClient) VerifyCodeHandler {
	return VerifyCodeHandler{grpcUser: grpcUser}
}

func (h VerifyCodeHandler) Handle(ctx context.Context, cmd VerifyCode) error {
	_, err := h.grpcUser.VerifyUser(ctx, &intrvproto.VerifyUserRequest{
		Code: cmd.VerificationCode,
	})
	if err != nil {
		return err
	}
	return nil
}

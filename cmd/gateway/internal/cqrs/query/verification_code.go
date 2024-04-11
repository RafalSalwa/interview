package query

import (
	"context"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type VerificationCodeHandler struct {
	authClient intrvproto.AuthServiceClient
}
type VerificationCode struct {
	Email string
}

func NewVerificationCodeHandler(authClient intrvproto.AuthServiceClient) VerificationCodeHandler {
	return VerificationCodeHandler{authClient: authClient}
}

func (h VerificationCodeHandler) Handle(ctx context.Context, email string) (models.UserResponse, error) {
	req := &intrvproto.VerificationCodeRequest{
		Email: email,
	}
	resp, err := h.authClient.GetVerificationKey(ctx, req)
	if err != nil {
		return models.UserResponse{}, err
	}
	u := models.UserResponse{
		VerificationCode: resp.GetCode(),
	}
	return u, nil
}

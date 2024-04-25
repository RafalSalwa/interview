package query

import (
	"context"

	"github.com/RafalSalwa/auth-api/pkg/models"
	intrvproto "github.com/RafalSalwa/auth-api/proto/grpc"
)

type GetUserByCodeHandler struct {
	grpcUser intrvproto.UserServiceClient
}

func NewGetUserByCodeHandler(userClient intrvproto.UserServiceClient) GetUserByCodeHandler {
	return GetUserByCodeHandler{grpcUser: userClient}
}

func (h GetUserByCodeHandler) Handle(ctx context.Context, vCode string) (models.UserResponse, error) {
	req := &intrvproto.VerificationCode{Code: vCode}

	pu, err := h.grpcUser.GetUserByCode(ctx, req)
	ur := models.UserResponse{}

	if err != nil {
		return ur, err
	}

	ur.FromProtoUserDetails(pu)

	return ur, nil
}

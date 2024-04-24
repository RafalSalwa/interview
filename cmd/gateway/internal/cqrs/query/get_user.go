package query

import (
	"context"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type GetUserHandler struct {
	grpcUser intrvproto.UserServiceClient
}

func NewGetUserHandler(userClient intrvproto.UserServiceClient) GetUserHandler {
	return GetUserHandler{grpcUser: userClient}
}

func (h GetUserHandler) Handle(ctx context.Context, user models.UserRequest) (models.UserResponse, error) {
	req := &intrvproto.GetUserRequest{
		Id:               user.Id,
		Email:            user.Email,
		VerificationCode: user.VerificationCode,
		Token:            user.AccessToken,
		RefreshToken:     user.RefreshToken}

	pu, err := h.grpcUser.GetUserById(ctx, req)
	ur := models.UserResponse{}

	if err != nil {
		return ur, err
	}

	ur.FromProtoUserDetails(pu)

	return ur, nil
}

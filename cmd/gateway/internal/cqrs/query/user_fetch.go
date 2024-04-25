package query

import (
	"context"

	"github.com/RafalSalwa/auth-api/pkg/models"
	intrvproto "github.com/RafalSalwa/auth-api/proto/grpc"
)

type FetchUserHandler struct {
	userClient intrvproto.UserServiceClient
}
type FetchUser struct {
	models.SignInUserRequest
}

func NewFetchUserHandler(userClient intrvproto.UserServiceClient) FetchUserHandler {
	return FetchUserHandler{userClient: userClient}
}

func (h FetchUserHandler) Handle(ctx context.Context, q FetchUser) (models.UserResponse, error) {
	credentials := &intrvproto.GetUserSignInRequest{
		Email:    q.Email,
		Password: q.Password,
	}

	resp, err := h.userClient.GetUser(ctx, credentials)
	if err != nil {
		return models.UserResponse{}, err
	}
	u := models.UserResponse{}
	u.FromProtoUserDetails(resp)
	return u, nil
}

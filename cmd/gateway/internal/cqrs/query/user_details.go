package query

import (
	"context"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type UserRequest struct {
	UserID int64
}

type UserDetailsHandler struct {
	grpcUser intrvproto.UserServiceClient
}

func NewUserDetailsHandler(userClient intrvproto.UserServiceClient) UserDetailsHandler {
	return UserDetailsHandler{grpcUser: userClient}
}

func (h UserDetailsHandler) Handle(ctx context.Context, query UserRequest) (*models.UserDBResponse, error) {
	req := &intrvproto.GetUserRequest{Id: query.UserID}
	pu, err := h.grpcUser.GetUserDetails(ctx, req)
	if err != nil {
		return nil, err
	}

	ur := &models.UserDBResponse{}
	err = ur.FromProtoUserDetails(pu)
	if err != nil {
		return nil, err
	}

	return ur, nil
}

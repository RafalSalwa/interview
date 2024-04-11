package command

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type SignInHandler struct {
	authClient intrvproto.AuthServiceClient
}

func NewSignInHandler(authClient intrvproto.AuthServiceClient) SignInHandler {
	return SignInHandler{authClient: authClient}
}

func (h SignInHandler) Handle(ctx context.Context, req models.SignInUserRequest) (*models.UserResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer("SignInUser").Start(ctx, "CQRS")
	defer span.End()

	credentials := &intrvproto.SignInUserInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}
	resp, err := h.authClient.SignInUser(ctx, credentials)
	if err != nil {
		return nil, err
	}
	u := &models.UserResponse{}
	u.FromProtoSignIn(resp)
	return u, nil
}

package command

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type SignInByCodeHandler struct {
	authClient intrvproto.AuthServiceClient
}

func NewSignInByCodeHandler(authClient intrvproto.AuthServiceClient) SignInByCodeHandler {
	return SignInByCodeHandler{authClient: authClient}
}

func (h SignInByCodeHandler) Handle(ctx context.Context, email, authCode string) (*models.UserResponse, error) {
	ctx, span := otel.GetTracerProvider().Tracer("SignInUser").Start(ctx, "CQRS")
	defer span.End()

	credentials := &intrvproto.SignInByCodeUserInput{
		Email:    email,
		AuthCode: authCode,
	}

	resp, err := h.authClient.SignInByCode(ctx, credentials)
	if err != nil {
		return nil, err
	}
	u := &models.UserResponse{}
	u.FromProtoSignIn(resp)
	return u, nil
}

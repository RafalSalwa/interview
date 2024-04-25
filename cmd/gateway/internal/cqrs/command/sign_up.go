package command

import (
	"context"

	"go.opentelemetry.io/otel"

	"github.com/RafalSalwa/auth-api/pkg/models"
	intrvproto "github.com/RafalSalwa/auth-api/proto/grpc"
)

type SignUpUser struct {
	User models.SignUpUserRequest
}

type SignUpHandler struct {
	authClient intrvproto.AuthServiceClient
}

func NewSignUpHandler(authClient intrvproto.AuthServiceClient) SignUpHandler {
	return SignUpHandler{authClient: authClient}
}

func (h SignUpHandler) Handle(ctx context.Context, req models.SignUpUserRequest) error {
	ctx, span := otel.GetTracerProvider().Tracer("CQRS").Start(ctx, "Command/SignUpUser")
	defer span.End()

	_, err := h.authClient.SignUpUser(ctx, &intrvproto.SignUpUserInput{
		Email:           req.Email,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
	})

	if err != nil {
		return err
	}
	return nil
}

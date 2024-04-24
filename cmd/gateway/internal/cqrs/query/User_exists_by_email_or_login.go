package query

import (
	"context"

	"go.opentelemetry.io/otel"

	intrvproto "github.com/RafalSalwa/interview-app-srv/proto/grpc"
)

type UserExistsHandler struct {
	userClient intrvproto.UserServiceClient
}

func NewUserExistsHandler(userClient intrvproto.UserServiceClient) UserExistsHandler {
	return UserExistsHandler{userClient: userClient}
}

func (h UserExistsHandler) Handle(ctx context.Context, email string) (bool, error) {
	ctx, span := otel.GetTracerProvider().Tracer("CQRS").Start(ctx, "Query/SignUpUser")
	defer span.End()

	req := &intrvproto.StringValue{Value: email}
	resp, err := h.userClient.CheckUserExists(ctx, req)
	if err != nil {
		return false, err
	}
	return resp.GetValue(), nil
}

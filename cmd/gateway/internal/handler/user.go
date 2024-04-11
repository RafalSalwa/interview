package handler

import (
	"net/http"
	"strconv"

	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/internal/cqrs"
	"github.com/RafalSalwa/interview-app-srv/pkg/hashing"
	"github.com/RafalSalwa/interview-app-srv/pkg/http/auth"
	"github.com/RafalSalwa/interview-app-srv/pkg/http/middlewares"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/models"
	"github.com/RafalSalwa/interview-app-srv/pkg/responses"
	"github.com/RafalSalwa/interview-app-srv/pkg/tracing"
	"github.com/RafalSalwa/interview-app-srv/pkg/validate"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"
)

type UserHandler interface {
	RouteRegisterer
	GetUserByID() HandlerFunc
	PasswordChange() HandlerFunc
}

type userHandler struct {
	cqrs   *cqrs.Application
	logger *logger.Logger
}

func NewUserHandler(cqrs *cqrs.Application, l *logger.Logger) UserHandler {
	return userHandler{cqrs, l}
}

func (uh userHandler) RegisterRoutes(r *mux.Router, cfg interface{}) {
	params := cfg.(auth.JWTConfig)
	s := r.PathPrefix("/user").Subrouter()
	s.Use(middlewares.ValidateJWTAccessToken(params))

	s.Methods(http.MethodGet).Path("").HandlerFunc(uh.GetUserByID())
	s.Methods(http.MethodPost).Path("/change_password").HandlerFunc(uh.PasswordChange())
}

func (uh userHandler) GetUserByID() HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("user-handler").Start(r.Context(), "GetUserByID")
		defer span.End()

		userID, err := strconv.ParseInt(r.Header.Get("x-user-id"), 10, 64)

		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			uh.logger.Error().Err(err).Msg("GetUserByID:header:getId")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		user, err := uh.cqrs.GetUser(ctx, models.UserRequest{Id: userID})
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			uh.logger.Error().Err(err).Msg("GetUserByID:grpc:getUser")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}
		responses.User(w, user)
	}
}

func (uh userHandler) PasswordChange() HandlerFunc {
	req := &models.ChangePasswordRequest{}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("user-handler").Start(r.Context(), "PasswordChange")
		defer span.End()

		userID, err := strconv.ParseInt(r.Header.Get("x-user-id"), 10, 64)
		if err != nil {
			tracing.RecordError(span, err)
			uh.logger.Error().Err(err).Msg("GetUserByID:header:getId")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		if err = validate.UserInput(r, &req); err != nil {
			tracing.RecordError(span, err)
			uh.logger.Error().Err(err).Msg("PasswordChange: decode")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		if err = hashing.Validate(req.Password, req.PasswordConfirm); err != nil {
			tracing.RecordError(span, err)
			uh.logger.Error().Err(err).Msg("PasswordChange:validateInputPasswords")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		_, err = uh.cqrs.GetUser(ctx, models.UserRequest{Id: userID})
		if err != nil {
			tracing.RecordError(span, err)
			uh.logger.Error().Err(err).Msg("PasswordChange:grpc:GetUserByID")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		err = uh.cqrs.ChangePassword(ctx, req)
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			uh.logger.Error().Err(err).Msg("PasswordChange:grpc:ChangePassword")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}

		responses.RespondOk(w)
	}
}

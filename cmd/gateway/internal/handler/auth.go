package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/RafalSalwa/auth-api/cmd/gateway/internal/cqrs"
	"github.com/RafalSalwa/auth-api/cmd/gateway/internal/cqrs/command"
	"github.com/RafalSalwa/auth-api/pkg/http/auth"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/models"
	"github.com/RafalSalwa/auth-api/pkg/responses"
	"github.com/RafalSalwa/auth-api/pkg/tracing"
	"github.com/RafalSalwa/auth-api/pkg/validate"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/status"
)

type AuthHandler interface {
	RouteRegisterer

	SignUpUser() http.HandlerFunc
	SignInUser() http.HandlerFunc

	Verify() http.HandlerFunc
	GetVerificationCode() http.HandlerFunc
}

type authHandler struct {
	cqrs   *cqrs.Application
	logger *logger.Logger
}

func (a authHandler) RegisterRoutes(r *mux.Router, cfg interface{}) {
	params := cfg.(auth.Auth)
	authorizer, _ := auth.NewAuthorizer(params)

	sr := r.PathPrefix("/auth/").Subrouter()

	sr.Methods(http.MethodPost).Path("/signup").HandlerFunc(authorizer.Middleware(a.SignUpUser()))
	sr.Methods(http.MethodPost).Path("/signin/{auth_code}").HandlerFunc(authorizer.Middleware(a.SignInUserByCode()))
	sr.Methods(http.MethodPost).Path("/signin").HandlerFunc(authorizer.Middleware(a.SignInUser()))

	sr.Methods(http.MethodGet).Path("/verify/{code}").HandlerFunc(authorizer.Middleware(a.Verify()))
	sr.Methods(http.MethodPost).Path("/code").HandlerFunc(authorizer.Middleware(a.GetVerificationCode()))
	sr.Methods(http.MethodGet).Path("/code/{code}").HandlerFunc(authorizer.Middleware(a.GetUserByCode()))
}

func NewAuthHandler(cqrs *cqrs.Application, l *logger.Logger) AuthHandler {
	return authHandler{cqrs, l}
}

func (a authHandler) SignInUser() http.HandlerFunc {
	reqUser := models.SignInUserRequest{}

	res := &models.UserResponse{}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("SignInUser").Start(r.Context(), "SignInUser Handler")
		defer span.End()

		if err := validate.UserInput(r, &reqUser); err != nil {
			a.logger.Error().Err(err).Msg("SignInUser: validate")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		var errQuery error
		res, errQuery = a.cqrs.SigninCommand(ctx, reqUser)

		if errQuery != nil {
			a.logger.Error().Err(errQuery).Msg("SignInUser: grpc signIn")

			if e, ok := status.FromError(errQuery); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.InternalServerError(w)
			return
		}
		responses.NewUserResponse(res, w)
	}
}

func (a authHandler) SignInUserByCode() http.HandlerFunc {
	var authCode string
	reqSignIn := models.VerificationCodeRequest{}
	res := &models.UserResponse{}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("SignInUser").Start(r.Context(), "SignInUser Handler")
		defer span.End()
		authCode = mux.Vars(r)["auth_code"]
		if authCode == "" {
			authCode = r.URL.Query().Get("auth_code")
			if authCode == "" {
				responses.RespondBadRequest(w, "code param is missing")
				return
			}
		}
		if err := validate.UserInput(r, &reqSignIn); err != nil {
			tracing.RecordError(span, err)
			a.logger.Error().Err(err).Msg("SignInUserByCode: decode")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		var errQuery error
		res, errQuery = a.cqrs.SigninByCodeCommand(ctx, reqSignIn.Email, authCode)

		if errQuery != nil {
			a.logger.Error().Err(errQuery).Msg("SignInUser: grpc signIn")

			if e, ok := status.FromError(errQuery); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.InternalServerError(w)
			return
		}
		responses.NewUserResponse(res, w)
	}
}

func (a authHandler) SignUpUser() http.HandlerFunc {
	var reqUser models.SignUpUserRequest

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("Handler").Start(r.Context(), "Handler/SignUpUser")
		defer span.End()

		if err := validate.UserInput(r, &reqUser); err != nil {
			tracing.RecordError(span, err)
			a.logger.Error().Err(err).Msg("SignUpUser: decode")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		err := a.cqrs.SignupUserCommand(ctx, reqUser)

		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			a.logger.Error().Err(err).Msg("SignUpUser:create")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}
		responses.RespondCreated(w)
	}
}

func (a authHandler) GetVerificationCode() http.HandlerFunc {
	reqSignIn := models.VerificationCodeRequest{}
	resp := models.UserResponse{}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("auth-handler").Start(r.Context(), "Handler SignUpUser")
		defer span.End()

		if err := validate.UserInput(r, &reqSignIn); err != nil {
			tracing.RecordError(span, err)
			a.logger.Error().Err(err).Msg("GetVerificationCode: validate")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		_, err := a.cqrs.GetUser(ctx, models.UserRequest{Email: reqSignIn.Email})
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			a.logger.Error().Err(err).Msg("GetVerificationCode: fetchUser")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}

		resp, err = a.cqrs.GetVerificationCode(ctx, reqSignIn.Email)
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			a.logger.Error().Err(err).Msg("GetVerificationCode: GetVerificationCode")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}
		responses.User(w, resp)
	}
}

func (a authHandler) GetUserByCode() http.HandlerFunc {
	var vCode string

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("user-handler").Start(r.Context(), "GetUserByCode")
		defer span.End()
		dt := time.Now()
		fmt.Println("Current date and time is: ", dt.String())
		vCode = mux.Vars(r)["code"]
		if vCode == "" {
			vCode = r.URL.Query().Get("code")
			if vCode == "" {
				responses.RespondBadRequest(w, "code param is missing")
				return
			}
		}

		user, err := a.cqrs.GetUserByCode(ctx, vCode)
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			a.logger.Error().Err(err).Msg("GetUserByID:header:getId")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			a.logger.Error().Err(err).Msg("GetUserByID:grpc:getUser")

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

func (a authHandler) Verify() http.HandlerFunc {
	var vCode string

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("auth-handler").Start(r.Context(), "Handler SignUpUser")
		defer span.End()

		vCode = mux.Vars(r)["code"]

		if vCode == "" {
			vCode = r.URL.Query().Get("code")
			if vCode == "" {
				responses.RespondBadRequest(w, "code param is missing")
				return
			}
		}

		err := a.cqrs.Commands.VerifyUserByCode.Handle(ctx, command.VerifyCode{VerificationCode: vCode})

		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			a.logger.Error().Err(err).Msg("Verify")

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

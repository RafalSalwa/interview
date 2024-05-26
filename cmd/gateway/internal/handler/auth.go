package handler

import (
	"net/http"

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

type (
	AuthHandler interface {
		RouteRegisterer

		SignUpUser() http.HandlerFunc
		SignInUser() http.HandlerFunc

		Verify() http.HandlerFunc
		GetVerificationCode() http.HandlerFunc
	}
	authHandler struct {
		application *cqrs.Application
		logger      *logger.Logger
	}
)

const verificationCodeParam = "code"

func (h authHandler) RegisterRoutes(r *mux.Router, cfg interface{}) {
	params := cfg.(auth.Auth)
	authorizer, _ := auth.NewAuthorizer(&params)

	sr := r.PathPrefix("/auth/").Subrouter()

	sr.Methods(http.MethodPost).Path("/signup").HandlerFunc(authorizer.Middleware(h.SignUpUser()))
	sr.Methods(http.MethodPost).Path("/signin/{auth_code}").HandlerFunc(authorizer.Middleware(h.SignInUserByCode()))
	sr.Methods(http.MethodPost).Path("/signin").HandlerFunc(authorizer.Middleware(h.SignInUser()))

	sr.Methods(http.MethodGet).Path("/verify/{code}").HandlerFunc(authorizer.Middleware(h.Verify()))
	sr.Methods(http.MethodPost).Path("/code").HandlerFunc(authorizer.Middleware(h.GetVerificationCode()))
	sr.Methods(http.MethodGet).Path("/code/{code}").HandlerFunc(authorizer.Middleware(h.GetUserByCode()))
}

func NewAuthHandler(application *cqrs.Application, l *logger.Logger) AuthHandler {
	return authHandler{application, l}
}

func (h authHandler) SignInUser() http.HandlerFunc {
	reqUser := models.SignInUserRequest{}

	res := &models.UserResponse{}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("SignInUser").Start(r.Context(), "SignInUser Handler")
		defer span.End()

		if err := validate.UserInput(r, &reqUser); err != nil {
			h.logger.Error().Err(err).Msg("SignInUser: validate")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		var errQuery error
		res, errQuery = h.application.SigninCommand(ctx, reqUser)

		if errQuery != nil {
			h.logger.Error().Err(errQuery).Msg("SignInUser: grpc signIn")

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

func (h authHandler) SignInUserByCode() http.HandlerFunc {
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
			h.logger.Error().Err(err).Msg("SignInUserByCode: decode")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		var errQuery error
		res, errQuery = h.application.SigninByCodeCommand(ctx, reqSignIn.Email, authCode)

		if errQuery != nil {
			h.logger.Error().Err(errQuery).Msg("SignInUser: grpc signIn")

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

func (h authHandler) SignUpUser() http.HandlerFunc {
	var reqUser models.SignUpUserRequest

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("Handler").Start(r.Context(), "AuthHandler/SignUpUser")
		defer span.End()

		if err := validate.UserInput(r, &reqUser); err != nil {
			tracing.RecordError(span, err)
			h.logger.Error().Err(err).Msg("SignUpUser: validate")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		err := h.application.SignupUserCommand(ctx, reqUser)

		if err != nil {
			tracing.RecordError(span, err)
			h.logger.Error().Err(err).Msg("SignUpUser:create")

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

func (h authHandler) GetVerificationCode() http.HandlerFunc {
	reqSignIn := models.VerificationCodeRequest{}
	resp := models.UserResponse{}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("auth-handler").Start(r.Context(), "Handler SignUpUser")
		defer span.End()

		if err := validate.UserInput(r, &reqSignIn); err != nil {
			tracing.RecordError(span, err)
			h.logger.Error().Err(err).Msg("GetVerificationCode: validate")

			responses.RespondBadRequest(w, err.Error())
			return
		}

		_, err := h.application.GetUser(ctx, models.UserRequest{Email: reqSignIn.Email})
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			h.logger.Error().Err(err).Msg("GetVerificationCode: fetchUser")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}

		resp, err = h.application.GetVerificationCode(ctx, reqSignIn.Email)
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			h.logger.Error().Err(err).Msg("GetVerificationCode: GetVerificationCode")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}
		responses.User(w, &resp)
	}
}

func (h authHandler) GetUserByCode() http.HandlerFunc {
	var vCode string

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("user-handler").Start(r.Context(), "GetUserByCode")
		defer span.End()

		vCode = mux.Vars(r)[verificationCodeParam]
		if vCode == "" {
			vCode = r.URL.Query().Get(verificationCodeParam)
			if vCode == "" {
				responses.RespondBadRequest(w, "code param is missing")
				return
			}
		}

		user, err := h.application.GetUserByCode(ctx, vCode)
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			h.logger.Error().Err(err).Msg("GetUserByID:header:getId")
			responses.RespondBadRequest(w, err.Error())
			return
		}

		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			h.logger.Error().Err(err).Msg("GetUserByID:grpc:getUser")

			if e, ok := status.FromError(err); ok {
				responses.FromGRPCError(e, w)
				return
			}
			responses.RespondBadRequest(w, err.Error())
			return
		}
		responses.User(w, &user)
	}
}

func (h authHandler) Verify() http.HandlerFunc {
	var vCode string

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.GetTracerProvider().Tracer("auth-handler").Start(r.Context(), "Handler SignUpUser")
		defer span.End()

		vCode = mux.Vars(r)[verificationCodeParam]

		if vCode == "" {
			vCode = r.URL.Query().Get(verificationCodeParam)
			if vCode == "" {
				responses.RespondBadRequest(w, "code param is missing")
				return
			}
		}

		err := h.application.Commands.VerifyUserByCode.Handle(ctx, command.VerifyCode{VerificationCode: vCode})

		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
			h.logger.Error().Err(err).Msg("Verify")

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

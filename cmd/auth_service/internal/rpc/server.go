package rpc

import (
	"github.com/RafalSalwa/auth-api/cmd/auth_service/internal/services"
	grpcconfig "github.com/RafalSalwa/auth-api/pkg/grpc"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	pb "github.com/RafalSalwa/auth-api/proto/grpc"
)

type Auth struct {
	pb.UnimplementedAuthServiceServer
	config      grpcconfig.Config
	logger      *logger.Logger
	authService services.AuthService
}

func NewGrpcAuthServer(config grpcconfig.Config, logger *logger.Logger, authService services.AuthService) (*Auth, error) {
	authServer := &Auth{
		config:      config,
		logger:      logger,
		authService: authService,
	}

	return authServer, nil
}

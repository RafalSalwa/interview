package rpc

import (
	"github.com/RafalSalwa/auth-api/cmd/user_service/internal/services"
	grpcconfig "github.com/RafalSalwa/auth-api/pkg/grpc"
	pb "github.com/RafalSalwa/auth-api/proto/grpc"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	config      grpcconfig.Config
	userService services.UserService
}

func NewGrpcUserServer(config grpcconfig.Config, userService services.UserService) (*UserServer, error) {
	userServer := &UserServer{
		config:      config,
		userService: userService,
	}

	return userServer, nil
}

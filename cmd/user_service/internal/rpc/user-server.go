package rpc

import (
	"github.com/RafalSalwa/interview-app-srv/cmd/user_service/internal/services"
	grpcconfig "github.com/RafalSalwa/interview-app-srv/pkg/grpc"
	pb "github.com/RafalSalwa/interview-app-srv/proto/grpc"
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

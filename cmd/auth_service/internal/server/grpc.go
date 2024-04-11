package server

import (
	"net"
	"time"

	"github.com/RafalSalwa/interview-app-srv/pkg/probes"
	"github.com/RafalSalwa/interview-app-srv/pkg/tracing"

	"github.com/RafalSalwa/interview-app-srv/cmd/auth_service/internal/rpc"
	"github.com/RafalSalwa/interview-app-srv/cmd/auth_service/internal/services"
	grpcconfig "github.com/RafalSalwa/interview-app-srv/pkg/grpc"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	pb "github.com/RafalSalwa/interview-app-srv/proto/grpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 1
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GRPC struct {
	pb.UnimplementedAuthServiceServer
	config      grpcconfig.Config
	probesCfg   *probes.Config
	logger      *logger.Logger
	authService services.AuthService
}

func NewGrpcServer(
	config grpcconfig.Config,
	log *logger.Logger,
	probesCfg *probes.Config,
	authService services.AuthService) *GRPC {
	srv := &GRPC{
		config:      config,
		logger:      log,
		probesCfg:   probesCfg,
		authService: authService,
	}
	return srv
}

func (s *GRPC) Run() {
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(
			grpcctxtags.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
			grpcrecovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			otelgrpc.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	authServer, err := rpc.NewGrpcAuthServer(s.config, s.logger, s.authService)
	if err != nil {
		s.logger.Error().Err(err).Msg("auth:server:new")
	}

	pb.RegisterAuthServiceServer(grpcServer, authServer)
	reflection.Register(grpcServer)
	tracing.RegisterMetricsEndpoint(s.probesCfg.Port)
	listener, err := net.Listen("tcp", s.config.Addr)
	if err != nil {
		s.logger.Error().Err(err).Msg("grpc:net:listen")
	}

	s.logger.Info().Msgf("Starting gRPC server on: %s", s.config.Addr)
	err = grpcServer.Serve(listener)
	if err != nil {
		s.logger.Error().Err(err).Msg("grpc:server:server")
	}
	grpcServer.GracefulStop()
}

package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RafalSalwa/interview-app-srv/cmd/user_service/config"
	"github.com/RafalSalwa/interview-app-srv/cmd/user_service/internal/services"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/tracing"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	log *logger.Logger
	cfg *config.Config
	v   *validator.Validate
}

func NewGRPC(cfg *config.Config, log *logger.Logger) *Server {
	return &Server{log: log, cfg: cfg, v: validator.New()}
}

func (srv *Server) Run() error {
	ctx, rejectContext := context.WithCancel(NewContextCancellableByOsSignals(context.Background()))

	userService := services.NewUserService(ctx, srv.cfg, srv.log)

	grpcServer, err := NewGrpcServer(srv.cfg.GRPC, srv.cfg.Probes, &userService)
	if err != nil {
		srv.log.Error().Err(err)
	}
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		grpcServer.Run(srv.log)
	}()

	if srv.cfg.Jaeger.Enable {
		if err := tracing.OTELGRPCProvider(srv.cfg.ServiceName, srv.cfg.Jaeger); err != nil {
			srv.log.Error().Err(err).Msg("server:jaeger:register")
		}
	}
	<-shutdown
	rejectContext()
	return nil
}

func NewContextCancellableByOsSignals(parent context.Context) context.Context {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	newCtx, cancel := context.WithCancel(parent)

	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			fmt.Println("Received Interrupt signal")
			cancel()
		case syscall.SIGTERM:
			fmt.Println("Received sigterm signal")
			cancel()
		}
	}()

	return newCtx
}

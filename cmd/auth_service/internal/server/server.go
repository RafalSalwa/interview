package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/RafalSalwa/auth-api/cmd/auth_service/config"
	"github.com/RafalSalwa/auth-api/cmd/auth_service/internal/services"
	"github.com/RafalSalwa/auth-api/pkg/logger"
	"github.com/RafalSalwa/auth-api/pkg/tracing"
)

type Server struct {
	log *logger.Logger
	cfg *config.Config
}

func NewGRPC(cfg *config.Config, log *logger.Logger) *Server {
	return &Server{log: log, cfg: cfg}
}

func (srv *Server) Run() error {
	ctx, rejectContext := context.WithCancel(NewContextCancellableByOsSignals(context.Background()))

	authService := services.NewAuthService(ctx, srv.cfg, srv.log)
	s := NewGrpcServer(srv.cfg.GRPC, srv.log, &srv.cfg.Probes, authService)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		s.Run()
	}()

	if err := tracing.OTELGRPCProvider(srv.cfg.ServiceName); err != nil {
		srv.log.Error().Err(err).Msg("server:jaeger:register")
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
			cancel()
		case syscall.SIGTERM:
			cancel()
		}
	}()

	return newCtx
}

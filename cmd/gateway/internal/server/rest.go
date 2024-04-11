package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/cmd/gateway/config"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
	"github.com/RafalSalwa/interview-app-srv/pkg/metrics"
	"github.com/RafalSalwa/interview-app-srv/pkg/tracing"
	"github.com/gorilla/mux"
)

type Server struct {
	srv *http.Server
	mtr *metrics.RequestCounter
	log *logger.Logger
	cfg *config.Config
}

func NewServer(cfg *config.Config, r *mux.Router, l *logger.Logger) *Server {
	tlsConf := new(tls.Config)
	s := &http.Server{
		Addr:         cfg.HTTP.Addr,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.TimeoutRead,
		WriteTimeout: cfg.HTTP.TimeoutWrite,
		IdleTimeout:  cfg.HTTP.TimeoutIdle,
		TLSConfig:    tlsConf,
	}

	return &Server{
		srv: s,
		log: l,
		mtr: metrics.NewRequestCounterMetrics(cfg.ServiceName),
		cfg: cfg,
	}
}

func (srv *Server) ServeHTTP() {
	go func() {
		srv.log.Info().Msgf("Starting server server on: %v [auth method: %s]", srv.srv.Addr, srv.cfg.Auth.AuthMethod)
		if srv.cfg.App.Env == "dev" {
			if err := srv.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				srv.log.Error().Err(err).Msg("server:Listen")
			}
		} else {
			if err := srv.srv.ListenAndServeTLS(
				"/etc/ssl/certs/server.crt",
				"/etc/ssl/private/server.key"); err != nil && !errors.Is(err, http.ErrServerClosed) {
				srv.log.Error().Err(err).Msg("server:ListenTLS")
			}
		}
	}()

	if srv.cfg.Jaeger.Enable {
		if err := tracing.OTELGRPCProvider(srv.cfg.ServiceName, srv.cfg.Jaeger); err != nil {
			srv.log.Error().Err(err).Msg("server:jaeger:register")
		}
	}
}
func (srv *Server) Shutdown() {
	closed := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), srv.srv.IdleTimeout)
	defer cancel()

	if err := srv.srv.Shutdown(ctx); err != nil {
		srv.log.Error().Err(err).Msg("server:shutdown")
	}

	close(closed)
}

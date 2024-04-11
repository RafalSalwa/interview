package main

import (
	"fmt"
	"os"

	"github.com/RafalSalwa/interview-app-srv/cmd/user_service/config"
	"github.com/RafalSalwa/interview-app-srv/cmd/user_service/internal/server"
	"github.com/RafalSalwa/interview-app-srv/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	l := logger.NewConsole()
	cfg, err := config.InitConfig()
	if err != nil {
		l.Error().Err(err).Msg("config Init")
		return err
	}

	srv := server.NewGRPC(cfg, l)

	if errSrv := srv.Run(); errSrv != nil {
		l.Error().Err(err).Msg("server run")
		return err
	}
	return nil
}

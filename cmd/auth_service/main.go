package main

import (
	"fmt"
	"os"

	"github.com/RafalSalwa/auth-api/cmd/auth_service/config"
	"github.com/RafalSalwa/auth-api/cmd/auth_service/internal/server"
	"github.com/RafalSalwa/auth-api/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
func run() error {
	cfg, err := config.InitConfig()
	if err != nil {
		return err
	}

	l := logger.NewConsole()

	srv := server.NewGRPC(cfg, l)

	if errSrv := srv.Run(); errSrv != nil {
		l.Error().Err(err).Msg("srv:run")
		return err
	}
	return nil
}

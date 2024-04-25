package workers

import (
	"context"

	"github.com/RafalSalwa/auth-api/pkg/logger"

	"github.com/RafalSalwa/auth-api/cmd/tester_service/config"
)

type WorkerRunner interface {
	Run()
}

func NewWorker(kind string) WorkerRunner {
	l := logger.NewConsole()
	cfg, err := config.InitConfig()
	if err != nil {
		l.Error().Err(err).Msg("config init err")
	}
	ctx := context.Background()

	switch kind {
	case "sequential":
		return NewSequential(ctx, cfg, l)
	case "ordered":
		return NewOrdered(ctx, cfg, l)
	case "daisy_chain":
		return NewDaisyChain(cfg, l)
	case "pool":
		return NewPool(cfg)
	}
	return nil
}

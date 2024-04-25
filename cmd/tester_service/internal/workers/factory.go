package workers

import (
	"context"
	"fmt"

	"github.com/RafalSalwa/auth-api/cmd/tester_service/config"
)

type WorkerRunner interface {
	Run()
}

func NewWorker(kind string) WorkerRunner {
	cfg, err := config.InitConfig()
	if err != nil {
		fmt.Println("config", err)
	}
	ctx := context.Background()

	switch kind {
	case "sequential":
		return NewSequential(ctx, cfg)
	case "daisy_chain":
		return NewDaisyChain(cfg)
	case "pool":
		return NewPool(cfg)
	}
	return nil
}

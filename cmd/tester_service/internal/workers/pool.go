package workers

import (
	"net/http"

	"github.com/RafalSalwa/interview-app-srv/cmd/tester_service/config"
)

type Pool struct {
	cfg    *config.Config
	client *http.Client
}

func NewPool(cfg *config.Config) WorkerRunner {
	return &Pool{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (s Pool) Run() {

}

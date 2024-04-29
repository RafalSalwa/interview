package main

import (
	"github.com/RafalSalwa/auth-api/cmd/tester_service/internal/workers"
)

type testUser struct {
	ValidationCode string
	Username       string
	Email          string
	Password       string
}

func main() {
	worker := workers.NewWorker("daisy_chain")
	worker.Run()
}

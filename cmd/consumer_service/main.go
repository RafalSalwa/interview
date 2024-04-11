package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RafalSalwa/interview-app-srv/cmd/consumer_service/config"
	amqpHandlers "github.com/RafalSalwa/interview-app-srv/cmd/consumer_service/internal/handlers"
	"github.com/RafalSalwa/interview-app-srv/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	con := rabbitmq.NewConnection(cfg.AMQP)
	con.Connect(ctx)
	defer con.Close(ctx)

	ctx, rejectContext := context.WithCancel(NewContextCancellableByOsSignals(context.Background()))
	var client rabbitmq.Client = rabbitmq.NewClient(con)

	client.SetHandler("customer_account_confirmation_requested", amqpHandlers.WrapHandleCustomerAccountRequestConfirmEmail)
	client.SetHandler("customer_account_confirmed", amqpHandlers.WrapHandleCustomerAccountConfirmedEmail)
	client.SetHandler("customer_password_reset_requested", amqpHandlers.WrapHandleCustomerPasswordResetRequestedEmail)
	client.SetHandler("customer_password_reset_succeeded", amqpHandlers.WrapHandleCustomerPasswordResetRequestedSuccessed)
	client.SetHandler("customer_logged_in", amqpHandlers.WrapHandleCustomerLoggedIn)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		args := amqp.Table{
			"x-dead-letter-exchange": "ex_dlx",
		}
		_ = client.HandleChannel(ctx, "interview", "rsinterview", args)
		rejectContext()
	}()

	wg.Wait()
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
			fmt.Println("Received SIGTERM signal")
			cancel()
		case syscall.SIGINT:
			fmt.Println("Received SIGINT signal")
			cancel()
		}
	}()
	return newCtx
}

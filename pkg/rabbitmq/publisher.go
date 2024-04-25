package rabbitmq

import (
	"context"
	"errors"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Connection *Connection
	Exchange   *Exchange
}

func NewPublisher(cfg Config) (*Publisher, error) {
	con := NewConnection(cfg)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	con.Connect(ctx)
	return &Publisher{Connection: con, Exchange: cfg.Exchange}, nil
}

func (p *Publisher) Disconnect() error {
	ctxDone, cancelDone := context.WithTimeout(context.Background(), time.Second*10)
	notifDone := p.Connection.Close(ctxDone)
	select {
	case <-notifDone:
	case <-ctxDone.Done():
		cancelDone()
		return errors.New("failed to close rabbitmq connection")
	}
	cancelDone()
	return nil
}

func (p *Publisher) Publish(ctx context.Context, mes amqp.Publishing) error {
	err := p.Connection.Channel.PublishWithContext(
		ctx,
		p.Exchange.Name,
		p.Exchange.RoutingKey,
		false,
		false,
		mes,
	)
	if err != nil {
		return err
	}

	return nil
}

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

const ConnTimeout = 10 * time.Second

func NewPublisher(ctx context.Context, cfg Config) (*Publisher, error) {
	con := NewConnection(cfg)
	ctx, cancel := context.WithTimeout(ctx, ConnTimeout)
	defer cancel()
	con.Connect(ctx)
	return &Publisher{Connection: con, Exchange: cfg.Exchange}, nil
}

func (p *Publisher) Disconnect() error {
	ctxDone, cancelDone := context.WithTimeout(context.Background(), ConnTimeout)
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

func (p *Publisher) Publish(ctx context.Context, mes *amqp.Publishing) error {
	err := p.Connection.Channel.PublishWithContext(
		ctx,
		p.Exchange.Name,
		p.Exchange.RoutingKey,
		false,
		false,
		*mes,
	)
	if err != nil {
		return err
	}

	return nil
}

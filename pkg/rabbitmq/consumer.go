package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	name     string
	channel  *amqp.Channel
	delivery <-chan amqp.Delivery
	handler  ConsumerHandler
}

type ConsumerHandler func(data []byte) (success bool)

var ConsumerCanceledByContextError = fmt.Errorf("consumer canceled by context")
var ConsumerMessageNotInitialized = fmt.Errorf("consumer received empty message")

func (l *Consumer) HandleSingleMessage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			{
				return ConsumerCanceledByContextError
			}
		case d := <-l.delivery:
			{
				if d.Acknowledger == nil {
					return ConsumerMessageNotInitialized
				}

				if l.handler(d.Body) {
					return d.Ack(false)
				}
				return d.Reject(true)
			}
		}
	}
}

func (l *Consumer) Handle(ctx context.Context) error {
	for {
		if err := l.HandleSingleMessage(ctx); err != nil {
			return err
		}
	}
}

func (l *Consumer) Close() error {
	return l.channel.Close()
}

func (l *Connection) CreateConsumer(channelName string, consumerName string, handler ConsumerHandler, args amqp.Table) (*Consumer, error) {
	amqpChannel, err := l.Connection.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := amqpChannel.QueueDeclare(channelName, true, false, false, false, args)
	if err != nil {
		_ = amqpChannel.Close()
		return nil, err
	}
	err = amqpChannel.Qos(1, 0, false)
	if err != nil {
		_ = amqpChannel.Close()
		return nil, err
	}

	delivery, err := amqpChannel.Consume(
		queue.Name,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		_ = amqpChannel.Close()
		return nil, err
	}
	return &Consumer{
		channel:  amqpChannel,
		name:     channelName,
		delivery: delivery,
		handler:  handler,
	}, nil
}

package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	SetHandler(eventName string, handler EventHandler)
	HandleChannel(ctx context.Context, channelName string, consumerName string, args amqp.Table) error
}

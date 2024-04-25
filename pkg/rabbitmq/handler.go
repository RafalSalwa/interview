package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/RafalSalwa/auth-api/pkg/logger"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IntrvClient struct {
	connection *Connection
	logger     *logger.Logger
	handlers   map[string]EventHandler
	debug      bool
}

func NewClient(connection *Connection, l *logger.Logger) *IntrvClient {
	return &IntrvClient{
		connection: connection,
		logger:     l,
		handlers:   make(map[string]EventHandler),
		debug:      false,
	}
}

func (c *IntrvClient) SetDebug(debug bool) {
	c.debug = debug
}

func (c *IntrvClient) SetHandler(eventName string, handler EventHandler) {
	c.handlers[eventName] = handler
}

func (c *IntrvClient) HandleChannel(ctx context.Context, channelName, consumerName string, args amqp.Table) error {
	consumer, err := c.connection.CreateConsumer(channelName, consumerName, c.handleEvent, args)
	if err != nil {
		return err
	}

	defer func(consumer *Consumer) {
		err = consumer.Close()
		if err != nil {
			c.logger.Error().Err(err)
		}
	}(consumer)
	return consumer.Handle(ctx)
}

func (c *IntrvClient) handleEvent(data []byte) (isSuccess bool) {
	// create new event and deserialize it
	event := Event{}

	err := json.Unmarshal(data, &event)
	if err != nil {
		c.logger.Error().Err(err)
	}
	if handler, ok := c.handlers[event.Name]; ok {
		if err = handler(event); err != nil {
			return false
		}
	} else {
		return true
	}
	return true
}

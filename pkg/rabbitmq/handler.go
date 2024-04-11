package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type IntrvClient struct {
	connection *Connection
	handlers   map[string]EventHandler
	debug      bool
}

func NewClient(connection *Connection) *IntrvClient {
	return &IntrvClient{
		connection: connection,
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

func (c *IntrvClient) HandleChannel(ctx context.Context, channelName string, consumerName string, args amqp.Table) error {
	consumer, err := c.connection.CreateConsumer(channelName, consumerName, c.handleEvent, args)
	if err != nil {
		return err
	}

	defer func(consumer *Consumer) error {
		err := consumer.Close()
		if err != nil {
			return err
		}
		return nil
	}(consumer)
	return consumer.Handle(ctx)
}

func (c *IntrvClient) handleEvent(data []byte) (isSuccess bool) {
	// create new event and deserialize it
	event := Event{}

	_ = json.Unmarshal(data, &event)

	if c.debug {
		str1 := string(data)
		fmt.Println("String =", str1)
		fmt.Printf("processing event %s\n %#v\n", event.Name, event)
	}

	if handler, ok := c.handlers[event.Name]; ok {
		if err := handler(event); err != nil {
			fmt.Printf("handler error for event %s\n %#v\n", event.Name, err)
			return false
		}
	} else {
		fmt.Printf("No handler for event %s\n %#v\n", event.Name, event)
		return true
	}
	// we have no handler for that type of event
	return true
}

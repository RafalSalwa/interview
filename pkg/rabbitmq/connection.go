package rabbitmq

import (
	"context"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection struct {
	Connection  *amqp.Connection
	Channel     *amqp.Channel
	credentials Credentials
	wg          *sync.WaitGroup
}

func NewConnection(cfg Config) *Connection {
	cred := Credentials{
		Username: cfg.Username,
		Password: cfg.Password,
		Addr:     cfg.Addr,
		VHost:    cfg.VHost,
		Exchange: cfg.Exchange,
	}

	return &Connection{
		credentials: cred,
		wg:          &sync.WaitGroup{},
	}
}

func (l *Connection) Connect(ctx context.Context) {
	for {
		notifyClose, err := l.connect()
		if err != nil {
			log.Printf("error connecting to rabbitmq: [%s]\n", err)
			time.Sleep(time.Second * 5)
			continue
		}
		log.Printf("[RabbitMQ] Connection established\n")
		// create queues, exchanges, consume from queue, etc
		select {
		case <-notifyClose:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (l *Connection) connect() (notify chan *amqp.Error, err error) {
	l.Connection, err = amqp.Dial(l.credentials.GetURL())
	if err != nil {
		return
	}
	l.Channel, err = l.Connection.Channel()
	if err != nil {
		return
	}
	notify = make(chan *amqp.Error)
	l.Connection.NotifyClose(notify)
	return
}

func (l *Connection) Close(ctx context.Context) (done chan struct{}) {
	done = make(chan struct{})

	go func() {
		defer close(done)
		select { // either waits for the messages to process or timeout from context
		case <-ctx.Done():
		}
		l.closeConnections()
	}()
	return
}

func (l *Connection) closeConnections() {
	var err error
	if l.Channel != nil {
		err = l.Channel.Close()
		if err != nil {
			log.Printf("Error closing consumer channel: [%s]\n", err)
		}
	}

	if l.Connection != nil {
		err = l.Connection.Close()
		if err != nil {
			log.Printf("Error closing connection: [%s]\n", err)
		}
	}
}

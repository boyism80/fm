package server

import (
	"sync"

	"github.com/boyism80/fm/common/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQDeliveryHandler func(msg amqp.Delivery) error

type RabbitMQConsumer struct {
	cfg config.RabbitMQEndpoint

	conn        *amqp.Connection
	ch          *amqp.Channel
	exchange    string
	routingKeys []string
	queueName   string
	consumerTag string
	handler     RabbitMQDeliveryHandler

	stopOnce sync.Once
	stopped  chan struct{}
}

func NewRabbitMQConsumer(cfg config.RabbitMQEndpoint, exchange, queueName, consumerTag string, handler RabbitMQDeliveryHandler, routingKeys ...string) *RabbitMQConsumer {
	keys := append([]string(nil), routingKeys...)
	return &RabbitMQConsumer{
		cfg:         cfg,
		exchange:    exchange,
		routingKeys: keys,
		queueName:   queueName,
		consumerTag: consumerTag,
		handler:     handler,
		stopped:     make(chan struct{}),
	}
}

func (c *RabbitMQConsumer) QueueName() string {
	return c.queueName
}

func (c *RabbitMQConsumer) Start() error {
	url := c.cfg.AMQPURL()
	if url == "" {
		return nil
	}
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return err
	}
	if err := ch.ExchangeDeclare(c.exchange, "direct", true, false, false, false, nil); err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}
	q, err := ch.QueueDeclare(c.queueName, true, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}
	for _, rk := range c.routingKeys {
		if rk == "" {
			continue
		}
		if err := ch.QueueBind(q.Name, rk, c.exchange, false, nil); err != nil {
			_ = ch.Close()
			_ = conn.Close()
			return err
		}
	}
	deliveries, err := ch.Consume(q.Name, c.consumerTag, false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		_ = conn.Close()
		return err
	}

	c.conn = conn
	c.ch = ch

	go c.consumeLoop(deliveries)
	return nil
}

func (c *RabbitMQConsumer) consumeLoop(deliveries <-chan amqp.Delivery) {
	defer close(c.stopped)
	for msg := range deliveries {
		if c.handler == nil {
			_ = msg.Ack(false)
			continue
		}
		if err := c.handler(msg); err != nil {
			_ = msg.Nack(false, false)
			continue
		}
		_ = msg.Ack(false)
	}
}

func (c *RabbitMQConsumer) Close() error {
	var closeErr error
	c.stopOnce.Do(func() {
		if c.ch != nil {
			if err := c.ch.Cancel(c.consumerTag, false); err != nil {
				closeErr = err
			}
		}
		if c.conn != nil {
			if err := c.conn.Close(); err != nil && closeErr == nil {
				closeErr = err
			}
		}
		<-c.stopped
	})
	return closeErr
}

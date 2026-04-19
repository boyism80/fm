package mq

import "github.com/boyism80/fm/common/config"

type JSONConsumer[S any] struct {
	Reg         S
	cfg         config.RabbitMQEndpoint
	exchange    string
	queueName   string
	consumerTag string
	routingKeys []string
	Dispatcher  *Dispatcher
	inner       *Consumer
}

func NewJSONConsumer[S any](cfg config.RabbitMQEndpoint, reg S, queueName, consumerTag string) *JSONConsumer[S] {
	return &JSONConsumer[S]{
		Reg:         reg,
		cfg:         cfg,
		exchange:    DirectExchange,
		queueName:   queueName,
		consumerTag: consumerTag,
		Dispatcher:  NewDispatcher(),
	}
}

func (j *JSONConsumer[S]) Route(routingKey string) *JSONConsumer[S] {
	if j != nil && routingKey != "" {
		j.routingKeys = append(j.routingKeys, routingKey)
	}
	return j
}

func (j *JSONConsumer[S]) Start() error {
	if j == nil {
		return nil
	}
	j.inner = NewConsumer(j.cfg, j.exchange, j.queueName, j.consumerTag, j.Dispatcher.AsDeliveryHandler(), j.routingKeys...)
	return j.inner.Start()
}

func (j *JSONConsumer[S]) Close() error {
	if j == nil || j.inner == nil {
		return nil
	}
	return j.inner.Close()
}

func (j *JSONConsumer[S]) QueueName() string {
	if j == nil {
		return ""
	}
	if j.inner != nil {
		return j.inner.QueueName()
	}
	return j.queueName
}

package mq

import (
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/config"
	c_actor "github.com/boyism80/fm/core/actor"
)

// Queue is one RabbitMQ queue a server consumes, with the handlers bound to its dispatcher.
type Queue struct {
	Name        string
	ConsumerTag string
	ActorName   string
	RoutingKeys []string
	Bind        func(d *Dispatcher)
}

// Consumers runs one RabbitActor per queue while the broker is enabled.
type Consumers struct {
	actors *c_actor.ActorRegistry
	broker config.RabbitMQEndpoint
	queues []Queue
	pids   []*actor.PID
}

func NewConsumers(actors *c_actor.ActorRegistry, broker config.RabbitMQEndpoint, queues []Queue) *Consumers {
	return &Consumers{actors: actors, broker: broker, queues: queues}
}

func (c *Consumers) Start() error {
	if !c.broker.Enabled() {
		return nil
	}
	for _, q := range c.queues {
		d := NewDispatcher()
		q.Bind(d)
		cfg := RabbitActorConfig{
			Root:        c.actors.Root(),
			Broker:      c.broker,
			Exchange:    DirectExchange,
			QueueName:   q.Name,
			ConsumerTag: q.ConsumerTag,
			RoutingKeys: q.RoutingKeys,
			Dispatcher:  d,
		}
		props := actor.PropsFromProducer(func() actor.Actor {
			return NewRabbitActor(cfg)
		})
		c.pids = append(c.pids, c.actors.GetOrCreateActor(q.ActorName, props))
		log.Printf("MQ RabbitActor running (%s)", q.Name)
	}
	return nil
}

func (c *Consumers) Stop() {
	for _, pid := range c.pids {
		c.actors.Root().Poison(pid)
	}
}

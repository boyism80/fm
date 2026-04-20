package mq

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type JSONHandler interface {
	EventType() string
	Handle(ctx actor.Context, msg amqp.Delivery, eventType string, payload json.RawMessage) error
}

type JSONHandlerConstructor[S any, H JSONHandler] interface {
	New(S) H
}

// Bind registers a JSON event handler on d for use with RabbitActor.Dispatcher.
func Bind[S any, C JSONHandlerConstructor[S, H], H JSONHandler](registry S, d *Dispatcher) {
	if d == nil {
		return
	}
	var constructor C
	h := constructor.New(registry)
	et := h.EventType()
	if et == "" {
		return
	}
	d.Register(et, func(ctx actor.Context, msg amqp.Delivery, eventType string, raw json.RawMessage) error {
		return h.Handle(ctx, msg, eventType, raw)
	})
	log.Printf("mq: registered handler for event_type=%s %T", et, h)
}

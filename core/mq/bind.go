package mq

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type JSONHandler interface {
	EventType() string
	Handle(ctx actor.Context, payload json.RawMessage) error
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
	d.Register(h)
	log.Printf("mq: registered handler for event_type=%s %T", et, h)
}

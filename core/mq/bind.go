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

func (d *Dispatcher) Bind[C JSONHandlerConstructor[S, H], S any, H JSONHandler](registry S) {
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

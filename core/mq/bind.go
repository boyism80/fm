package mq

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type JSONHandler interface {
	EventType() string
	Handle(msg amqp.Delivery, eventType string, payload json.RawMessage) error
}

type JSONHandlerConstructor[S any, H JSONHandler] interface {
	New(S) H
}

func BindOn[S any, C JSONHandlerConstructor[S, H], H JSONHandler](j *JSONConsumer[S]) {
	if j == nil {
		return
	}
	Bind[S, C, H](j.Reg, j.Dispatcher)
}

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
	d.Register(et, func(msg amqp.Delivery, eventType string, raw json.RawMessage) error {
		return h.Handle(msg, eventType, raw)
	})
	log.Printf("mq: registered handler for event_type=%s %T", et, h)
}

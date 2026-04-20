package mq

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Handler processes one JSON event. ctx is the RabbitActor mailbox context when dispatched from AMQP.
type Handler func(ctx actor.Context, msg amqp.Delivery, eventType string, payload json.RawMessage) error

type Dispatcher struct {
	mu      sync.RWMutex
	byEvent map[string]Handler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		byEvent: make(map[string]Handler),
	}
}

func (d *Dispatcher) Register(eventType string, h Handler) {
	if d == nil || eventType == "" || h == nil {
		return
	}
	d.mu.Lock()
	d.byEvent[eventType] = h
	d.mu.Unlock()
}

// Dispatch invokes the handler for the JSON body (RabbitActor only).
func (d *Dispatcher) Dispatch(ctx actor.Context, body []byte) error {
	if d == nil {
		return nil
	}
	var head struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(body, &head); err != nil {
		log.Printf("mq: invalid JSON: %v", err)
		return nil
	}
	if head.EventType == "" {
		log.Printf("mq: missing event_type")
		return nil
	}

	d.mu.RLock()
	h, ok := d.byEvent[head.EventType]
	d.mu.RUnlock()
	if !ok {
		log.Printf("mq: unknown event_type=%s", head.EventType)
		return nil
	}
	raw := json.RawMessage(body)
	var zero amqp.Delivery
	if err := h(ctx, zero, head.EventType, raw); err != nil {
		log.Printf("mq: handler event_type=%s: %v", head.EventType, err)
	}
	return nil
}

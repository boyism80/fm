package mq

import (
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Handler func(msg amqp.Delivery, eventType string, payload json.RawMessage) error

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

func (d *Dispatcher) HandleDelivery(msg amqp.Delivery) error {
	if d == nil {
		return nil
	}
	var head struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(msg.Body, &head); err != nil {
		log.Printf("mq: invalid JSON (ack): %v", err)
		return nil
	}
	if head.EventType == "" {
		log.Printf("mq: missing event_type (ack)")
		return nil
	}

	d.mu.RLock()
	h, ok := d.byEvent[head.EventType]
	d.mu.RUnlock()
	if !ok {
		log.Printf("mq: unknown event_type=%s (ack)", head.EventType)
		return nil
	}
	raw := json.RawMessage(msg.Body)
	if err := h(msg, head.EventType, raw); err != nil {
		log.Printf("mq: handler event_type=%s (ack): %v", head.EventType, err)
	}
	return nil
}

func (d *Dispatcher) AsDeliveryHandler() DeliveryHandler {
	return d.HandleDelivery
}

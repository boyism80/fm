package mq

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

type Dispatcher struct {
	mu      sync.RWMutex
	byEvent map[string]JSONHandler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		byEvent: make(map[string]JSONHandler),
	}
}

func (d *Dispatcher) Register(h JSONHandler) {
	d.mu.Lock()
	d.byEvent[h.EventType()] = h
	d.mu.Unlock()
}

// Dispatch invokes the handler for the JSON body (RabbitActor only).
func (d *Dispatcher) Dispatch(ctx actor.Context, body []byte) {
	var head struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(body, &head); err != nil {
		log.Printf("mq: invalid JSON: %v", err)
		return
	}
	if head.EventType == "" {
		log.Printf("mq: missing event_type")
		return
	}

	d.mu.RLock()
	h, ok := d.byEvent[head.EventType]
	d.mu.RUnlock()
	if !ok {
		log.Printf("mq: unknown event_type=%s", head.EventType)
		return
	}
	if err := h.Handle(ctx, json.RawMessage(body)); err != nil {
		log.Printf("mq: handler event_type=%s: %v", head.EventType, err)
	}
}

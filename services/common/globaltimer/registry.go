package globaltimer

import (
	"context"
	"log"
	"sync"
	"time"
)

type Handler interface {
	GetName() string
	GetInterval() time.Duration
	GetInitialDelay() time.Duration
	Handle(ctx context.Context) error
}

type Registry struct {
	mu       sync.Mutex
	handlers map[string]Handler
}

func NewRegistry() *Registry {
	return &Registry{
		handlers: make(map[string]Handler),
	}
}

func (r *Registry) Register(h Handler) {
	if h == nil {
		return
	}
	r.mu.Lock()
	r.handlers[h.GetName()] = h
	r.mu.Unlock()
}

func RegisterTimer[H interface {
	Handler
	New() H
}](registry *Registry) {
	var zero H
	handler := zero.New()
	registry.Register(handler)
}

func (r *Registry) Start(ctx context.Context) {
	r.mu.Lock()
	list := make([]Handler, 0, len(r.handlers))
	for _, h := range r.handlers {
		list = append(list, h)
	}
	r.mu.Unlock()
	for _, h := range list {
		go runLoop(ctx, h)
	}
}

func runLoop(ctx context.Context, h Handler) {
	init := h.GetInitialDelay()
	if init < 0 {
		init = 0
	}
	iv := h.GetInterval()
	if iv <= 0 {
		return
	}
	if init > 0 {
		select {
		case <-ctx.Done():
			return
		case <-time.After(init):
		}
	}
	for {
		if err := h.Handle(ctx); err != nil {
			log.Printf("global timer %s: %v", h.GetName(), err)
		}
		select {
		case <-ctx.Done():
			return
		case <-time.After(iv):
		}
	}
}

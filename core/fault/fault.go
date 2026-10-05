package fault

import (
	"context"
	"strings"
	"sync"
	"time"
)

type Mode int

const (
	None Mode = iota
	Unreachable
	Lost
	Timeout
)

func (m Mode) String() string {
	switch m {
	case Unreachable:
		return "unreachable"
	case Lost:
		return "lost"
	case Timeout:
		return "timeout"
	default:
		return "off"
	}
}

type Injector struct {
	mutex  sync.RWMutex
	delay  time.Duration
	delays map[string]time.Duration
	modes  map[string]Mode
}

func NewInjector() *Injector {
	return &Injector{
		delays: make(map[string]time.Duration),
		modes:  make(map[string]Mode),
	}
}

func (in *Injector) Delay() time.Duration {
	in.mutex.RLock()
	defer in.mutex.RUnlock()

	return in.delay
}

func (in *Injector) SetDelay(delay time.Duration) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	in.delay = delay
}

func (in *Injector) KeyDelay(key string) time.Duration {
	in.mutex.RLock()
	defer in.mutex.RUnlock()

	return in.delays[key]
}

func (in *Injector) SetKeyDelay(key string, delay time.Duration) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if delay <= 0 {
		delete(in.delays, key)
		return
	}
	in.delays[key] = delay
}

func (in *Injector) Mode(key string) Mode {
	in.mutex.RLock()
	defer in.mutex.RUnlock()

	return in.modes[key]
}

func (in *Injector) SetMode(key string, mode Mode) {
	in.mutex.Lock()
	defer in.mutex.Unlock()

	if mode == None {
		delete(in.modes, key)
		return
	}
	in.modes[key] = mode
}

func (in *Injector) find(key string) (time.Duration, Mode) {
	in.mutex.RLock()
	defer in.mutex.RUnlock()

	delay := in.delay
	matched := -1
	for prefix, d := range in.delays {
		if len(prefix) > matched && strings.HasPrefix(key, prefix) {
			delay = d
			matched = len(prefix)
		}
	}

	mode := None
	matched = -1
	for prefix, m := range in.modes {
		if len(prefix) > matched && strings.HasPrefix(key, prefix) {
			mode = m
			matched = len(prefix)
		}
	}
	return delay, mode
}

func (in *Injector) wait(ctx context.Context, delay time.Duration) error {
	if delay <= 0 {
		return nil
	}

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

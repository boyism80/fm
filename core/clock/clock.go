package clock

import (
	"sync"
	"time"
)

var (
	mu     sync.RWMutex
	offset time.Duration
)

func Now() time.Time {
	mu.RLock()
	o := offset
	mu.RUnlock()
	return time.Now().In(time.Local).Add(o)
}

func Offset() time.Duration {
	mu.RLock()
	defer mu.RUnlock()
	return offset
}

func setOffset(d time.Duration) {
	mu.Lock()
	offset = d
	mu.Unlock()
}

func Reset() {
	setOffset(0)
}

func SetAbsolute(target time.Time) error {
	parsed, err := normalizeLocal(target)
	if err != nil {
		return err
	}
	wall := time.Now().In(time.Local)
	setOffset(parsed.Sub(wall))
	return nil
}

func AddOffset(d time.Duration) {
	mu.Lock()
	offset += d
	mu.Unlock()
}

func normalizeLocal(t time.Time) (time.Time, error) {
	if t.IsZero() {
		return time.Time{}, ErrInvalidDateTime
	}
	return t.In(time.Local), nil
}

package entity

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var ErrMapClosed = errors.New("map is closed")

const InstanceLease = 60 * time.Second

// instanceState is nil on shared maps.
type instanceState struct {
	mu      sync.Mutex
	refs    int
	closing bool
	closed  bool
	lease   *MapRef
}

// MapRef keeps an instance map open: characters on their way in and owners such as state machines each hold one.
type MapRef struct {
	m        *Map
	released atomic.Bool
}

func (r *MapRef) Release() {
	if r == nil || r.m == nil {
		return
	}
	if r.released.Swap(true) {
		return
	}
	r.m.release()
}

// Reserve fails once the map is closing, so a warp either gets a ref before the map can close or never starts.
// Shared maps never close and return a nil ref.
func (m *Map) Reserve() (*MapRef, error) {
	if !m.IsInstance() {
		return nil, nil
	}
	m.instance.mu.Lock()
	defer m.instance.mu.Unlock()
	if m.instance.closing || m.instance.closed {
		return nil, ErrMapClosed
	}
	m.instance.refs++
	return &MapRef{m: m}, nil
}

// StartLease keeps a new instance open until its first character enters. If nobody enters within InstanceLease, the map closes.
// Once anyone has entered the lease is gone, so the map closes the moment the last character leaves.
func (m *Map) StartLease() error {
	ref, err := m.Reserve()
	if err != nil {
		return err
	}
	if ref == nil {
		return nil
	}
	m.instance.mu.Lock()
	m.instance.lease = ref
	m.instance.mu.Unlock()

	time.AfterFunc(InstanceLease, m.endLease)
	return nil
}

func (m *Map) endLease() {
	m.instance.mu.Lock()
	lease := m.instance.lease
	m.instance.lease = nil
	m.instance.mu.Unlock()

	lease.Release()
}

// Close stops new refs; the map is removed when the last character leaves and the last ref is released.
func (m *Map) Close() {
	if !m.IsInstance() {
		return
	}
	m.instance.mu.Lock()
	defer m.instance.mu.Unlock()
	m.instance.closing = true
}

func (m *Map) Closing() bool {
	if !m.IsInstance() {
		return false
	}
	m.instance.mu.Lock()
	defer m.instance.mu.Unlock()
	return m.instance.closing || m.instance.closed
}

func (m *Map) enter() error {
	if !m.IsInstance() {
		return nil
	}
	m.instance.mu.Lock()
	if m.instance.closed {
		m.instance.mu.Unlock()
		return ErrMapClosed
	}
	m.instance.refs++
	lease := m.instance.lease
	m.instance.lease = nil
	m.instance.mu.Unlock()

	// The character's own ref was counted above, so dropping the lease cannot close the map.
	lease.Release()
	return nil
}

func (m *Map) release() {
	if !m.IsInstance() {
		return
	}
	m.instance.mu.Lock()
	m.instance.refs--
	closed := m.instance.refs == 0 && !m.instance.closed
	if closed {
		m.instance.closed = true
	}
	m.instance.mu.Unlock()

	if closed {
		m.GameWorld.GetMapSystem().RemoveInstanceMap(m.GetMapID())
	}
}

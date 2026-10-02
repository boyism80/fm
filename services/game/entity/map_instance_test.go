package entity

import (
	"errors"
	"sync"
	"testing"
)

type removeRecorder struct {
	MapSystem
	mu      sync.Mutex
	removed []uint32
}

func (r *removeRecorder) RemoveInstanceMap(instanceKey uint32) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.removed = append(r.removed, instanceKey)
	return nil
}

type recorderWorld struct {
	GameWorld
	maps *removeRecorder
}

func (w recorderWorld) GetMapSystem() MapSystem { return w.maps }

func newTestInstance(id uint32) (*Map, *removeRecorder) {
	rec := &removeRecorder{}
	return &Map{id: id, instance: &instanceState{}, GameWorld: recorderWorld{maps: rec}}, rec
}

func TestSharedMapNeverCloses(t *testing.T) {
	m := &Map{id: 1}
	ref, err := m.Reserve()
	if err != nil || ref != nil {
		t.Fatalf("shared Reserve = %v, %v", ref, err)
	}
	ref.Release()
	m.Close()
	if m.Closing() || m.Closed() {
		t.Fatal("shared map must never close")
	}
}

func TestInstanceClosesOnLastRelease(t *testing.T) {
	m, rec := newTestInstance(7)
	creator, err := m.Reserve()
	if err != nil {
		t.Fatal(err)
	}
	ticket, err := m.Reserve()
	if err != nil {
		t.Fatal(err)
	}

	if err := m.enter(); err != nil {
		t.Fatal(err)
	}
	ticket.Release()
	creator.Release()
	if m.Closed() || len(rec.removed) != 0 {
		t.Fatal("a character on the map keeps it open")
	}

	m.release()
	if m.Closed() == false || len(rec.removed) != 1 || rec.removed[0] != 7 {
		t.Fatalf("closed=%v removed=%v", m.Closed(), rec.removed)
	}
	if _, err := m.Reserve(); errors.Is(err, ErrMapClosed) == false {
		t.Fatalf("Reserve after close = %v", err)
	}
	if err := m.enter(); errors.Is(err, ErrMapClosed) == false {
		t.Fatalf("enter after close = %v", err)
	}
}

func TestDoubleReleaseCountsOnce(t *testing.T) {
	m, rec := newTestInstance(3)
	a, _ := m.Reserve()
	b, _ := m.Reserve()
	a.Release()
	a.Release()
	if m.Closed() || len(rec.removed) != 0 {
		t.Fatal("double release must not drop another holder's ref")
	}
	b.Release()
	if len(rec.removed) != 1 {
		t.Fatalf("removed=%v", rec.removed)
	}
}

func TestCloseRefusesNewRefsButKeepsInFlight(t *testing.T) {
	m, rec := newTestInstance(5)
	owner, _ := m.Reserve()
	ticket, _ := m.Reserve()

	m.Close()
	if _, err := m.Reserve(); errors.Is(err, ErrMapClosed) == false {
		t.Fatalf("Reserve while closing = %v", err)
	}
	if err := m.enter(); err != nil {
		t.Fatalf("a ticket reserved before Close must still enter: %v", err)
	}
	ticket.Release()
	owner.Release()
	if m.Closed() {
		t.Fatal("character still on the map")
	}
	m.release()
	if m.Closed() == false || len(rec.removed) != 1 {
		t.Fatalf("closed=%v removed=%v", m.Closed(), rec.removed)
	}
}

func TestLeaseHoldsUntilFirstEntry(t *testing.T) {
	m, rec := newTestInstance(11)
	if err := m.StartLease(); err != nil {
		t.Fatal(err)
	}
	if m.Closed() {
		t.Fatal("the lease keeps a never-entered instance open")
	}

	if err := m.enter(); err != nil {
		t.Fatal(err)
	}
	m.release()
	if m.Closed() == false || len(rec.removed) != 1 {
		t.Fatalf("first entry voids the lease, so the last leave closes at once: closed=%v removed=%v", m.Closed(), rec.removed)
	}
}

func TestLeaseExpiryClosesNeverEnteredInstance(t *testing.T) {
	m, rec := newTestInstance(12)
	if err := m.StartLease(); err != nil {
		t.Fatal(err)
	}
	m.endLease()
	if m.Closed() == false || len(rec.removed) != 1 {
		t.Fatalf("closed=%v removed=%v", m.Closed(), rec.removed)
	}
}

func TestLeaseExpiryAfterEntryIsNoop(t *testing.T) {
	m, rec := newTestInstance(13)
	if err := m.StartLease(); err != nil {
		t.Fatal(err)
	}
	if err := m.enter(); err != nil {
		t.Fatal(err)
	}
	m.endLease()
	if m.Closed() || len(rec.removed) != 0 {
		t.Fatal("the timer must mean nothing once someone has entered")
	}
	m.release()
	if len(rec.removed) != 1 {
		t.Fatalf("removed=%v", rec.removed)
	}
}

func TestConcurrentReleaseRemovesOnce(t *testing.T) {
	m, rec := newTestInstance(9)
	refs := make([]*MapRef, 64)
	for i := range refs {
		refs[i], _ = m.Reserve()
	}
	var wg sync.WaitGroup
	for _, ref := range refs {
		wg.Add(2)
		go func(r *MapRef) { defer wg.Done(); r.Release() }(ref)
		go func(r *MapRef) { defer wg.Done(); r.Release() }(ref)
	}
	wg.Wait()
	if len(rec.removed) != 1 {
		t.Fatalf("removed=%v", rec.removed)
	}
}

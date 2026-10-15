package game

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/route"
	lua "github.com/yuin/gopher-lua"
)

const (
	EnvelopeTimeout    = 5 * time.Second
	envelopeRetryDelay = 50 * time.Millisecond
)

// DispatchSystem knows which map actor owns each character on this channel and carries messages to it.
type DispatchSystem struct {
	gw           GameWorld
	mu           sync.RWMutex
	byCID        map[uint32]*dispatchEntry
	byName       map[string]uint32
	envelopeMu   sync.Mutex
	envelopes    map[uint64]*route.Envelope
	nextEnvelope atomic.Uint64
}

type dispatchEntry struct {
	name      string
	mapPID    *actor.PID
	envelopes []*route.Envelope
}

func NewDispatchSystem(gw GameWorld) *DispatchSystem {
	return &DispatchSystem{
		gw:        gw,
		byCID:     make(map[uint32]*dispatchEntry),
		byName:    make(map[string]uint32),
		envelopes: make(map[uint64]*route.Envelope),
	}
}

func (s *DispatchSystem) Start() error { return nil }
func (s *DispatchSystem) Stop()        {}

func nameKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func (s *DispatchSystem) RegisterCharacter(characterID uint32, name string) error {
	if characterID == 0 {
		return fmt.Errorf("dispatch: invalid character id 0")
	}
	key := nameKey(name)
	if key == "" {
		return fmt.Errorf("dispatch: empty character name")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.byCID[characterID]; exists {
		return fmt.Errorf("dispatch: character %d already registered", characterID)
	}
	if existingID, taken := s.byName[key]; taken && existingID != characterID {
		return fmt.Errorf("dispatch: name %q already in use", key)
	}
	s.byCID[characterID] = &dispatchEntry{name: key}
	s.byName[key] = characterID
	return nil
}

func (s *DispatchSystem) UnregisterCharacter(characterID uint32) {
	if characterID == 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ent, ok := s.byCID[characterID]
	if !ok {
		return
	}
	delete(s.byName, ent.name)
	delete(s.byCID, characterID)
}

func (s *DispatchSystem) SetMapPID(characterID uint32, pid *actor.PID) error {
	if characterID == 0 {
		return fmt.Errorf("dispatch: invalid character id 0")
	}
	s.mu.Lock()
	ent, ok := s.byCID[characterID]
	if !ok {
		s.mu.Unlock()
		return fmt.Errorf("dispatch: character %d not registered", characterID)
	}
	ent.mapPID = pid
	s.mu.Unlock()

	if pid != nil {
		s.flushEnvelopes(characterID, pid)
	}
	return nil
}

func (s *DispatchSystem) Exists(characterID uint32) bool {
	if characterID == 0 {
		return false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.byCID[characterID]
	return ok
}

func (s *DispatchSystem) GetMapPID(characterID uint32) (pid *actor.PID, ok bool) {
	if characterID == 0 {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	ent, exists := s.byCID[characterID]
	if !exists {
		return nil, false
	}
	return ent.mapPID, true
}

func (s *DispatchSystem) FindCharacterID(name string) (uint32, bool) {
	key := nameKey(name)
	if key == "" {
		return 0, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.byName[key]
	return id, ok
}

func (s *DispatchSystem) SendTo(characterID uint32, msg interface{}) {
	s.SendToCharacter(nil, characterID, msg)
}

func (s *DispatchSystem) Call(characterID uint32, run func(ctx actor.Context)) {
	s.SendTo(characterID, &MapCall{Run: func(ctx actor.Context, _ []*Map) []lua.LValue {
		run(ctx)
		return nil
	}})
}

func (s *DispatchSystem) CallCharacter(characterID uint32, run func(ctx actor.Context, ch *Character)) {
	s.SendTo(characterID, &MapCall{Run: func(ctx actor.Context, maps []*Map) []lua.LValue {
		for _, m := range maps {
			if ch := m.GetPlayer(characterID); ch != nil {
				run(ctx, ch)
				return nil
			}
		}
		run(ctx, nil)
		return nil
	}})
}

func (s *DispatchSystem) SendToCharacter(caller *actor.PID, characterID uint32, inner interface{}) {
	if characterID == 0 || inner == nil {
		return
	}
	corr := s.nextEnvelope.Add(1)
	if !s.Exists(characterID) {
		s.tellEnvelopeResult(caller, corr, false, "not_registered")
		return
	}

	d := &route.Envelope{
		CorrelationID:    corr,
		CharacterID:      characterID,
		DeadlineUnixNano: clock.Now().Add(EnvelopeTimeout).UnixNano(),
		Caller:           caller,
		Inner:            inner,
	}
	s.registerEnvelope(d)

	mapPID, ok := s.GetMapPID(characterID)
	if !ok {
		s.failEnvelope(corr, "not_registered")
		return
	}
	if mapPID != nil {
		s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(mapPID, d)
		return
	}
	if err := s.queueEnvelope(characterID, d); err != nil {
		s.failEnvelope(corr, "enqueue_failed")
	}
}

func (s *DispatchSystem) Resend(d *route.Envelope) {
	if d == nil {
		return
	}
	time.AfterFunc(envelopeRetryDelay, func() {
		s.retryEnvelope(d)
	})
}

func (s *DispatchSystem) CompleteEnvelope(correlationID uint64) {
	if correlationID == 0 {
		return
	}
	s.envelopeMu.Lock()
	delete(s.envelopes, correlationID)
	s.envelopeMu.Unlock()
}

func (s *DispatchSystem) registerEnvelope(d *route.Envelope) {
	s.envelopeMu.Lock()
	s.envelopes[d.CorrelationID] = d
	s.envelopeMu.Unlock()

	corr := d.CorrelationID
	time.AfterFunc(EnvelopeTimeout, func() {
		s.failEnvelope(corr, "timeout")
	})
}

func (s *DispatchSystem) failEnvelope(correlation uint64, reason string) {
	if correlation == 0 {
		return
	}
	s.envelopeMu.Lock()
	d, ok := s.envelopes[correlation]
	if ok {
		delete(s.envelopes, correlation)
	}
	s.envelopeMu.Unlock()
	if !ok {
		return
	}
	s.removeEnvelope(d.CharacterID, correlation)
	s.tellEnvelopeResult(d.Caller, correlation, false, reason)
}

func (s *DispatchSystem) retryEnvelope(d *route.Envelope) {
	if clock.Now().UnixNano() > d.DeadlineUnixNano {
		s.failEnvelope(d.CorrelationID, "timeout")
		return
	}
	mapPID, ok := s.GetMapPID(d.CharacterID)
	if !ok {
		s.failEnvelope(d.CorrelationID, "gone")
		return
	}
	if mapPID == nil {
		time.AfterFunc(envelopeRetryDelay, func() {
			s.retryEnvelope(d)
		})
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(mapPID, d)
}

func (s *DispatchSystem) flushEnvelopes(characterID uint32, mapPID *actor.PID) {
	for _, d := range s.takeEnvelopes(characterID) {
		if d == nil {
			continue
		}
		if clock.Now().UnixNano() > d.DeadlineUnixNano {
			s.failEnvelope(d.CorrelationID, "timeout")
			continue
		}
		s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(mapPID, d)
	}
}

func (s *DispatchSystem) AbandonEnvelopes(characterID uint32) {
	if characterID == 0 {
		return
	}
	for _, d := range s.takeEnvelopes(characterID) {
		if d != nil {
			s.failEnvelope(d.CorrelationID, "gone")
		}
	}
	s.envelopeMu.Lock()
	var rest []*route.Envelope
	for corr, d := range s.envelopes {
		if d != nil && d.CharacterID == characterID {
			rest = append(rest, d)
			delete(s.envelopes, corr)
		}
	}
	s.envelopeMu.Unlock()
	for _, d := range rest {
		s.tellEnvelopeResult(d.Caller, d.CorrelationID, false, "gone")
	}
}

func (s *DispatchSystem) tellEnvelopeResult(caller *actor.PID, correlation uint64, ok bool, reason string) {
	if caller == nil || correlation == 0 {
		return
	}
	s.gw.Systems().Get[*c_actor.ActorRegistry]().Root().Send(caller, &route.Result{
		CorrelationID: correlation,
		OK:            ok,
		Reason:        reason,
	})
}

func (s *DispatchSystem) queueEnvelope(characterID uint32, d *route.Envelope) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ent, ok := s.byCID[characterID]
	if !ok {
		return fmt.Errorf("dispatch: character %d not registered", characterID)
	}
	ent.envelopes = append(ent.envelopes, d)
	return nil
}

func (s *DispatchSystem) takeEnvelopes(characterID uint32) []*route.Envelope {
	s.mu.Lock()
	defer s.mu.Unlock()
	ent, ok := s.byCID[characterID]
	if !ok || len(ent.envelopes) == 0 {
		return nil
	}
	buf := ent.envelopes
	ent.envelopes = nil
	return buf
}

func (s *DispatchSystem) removeEnvelope(characterID uint32, correlationID uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ent, ok := s.byCID[characterID]
	if !ok || len(ent.envelopes) == 0 {
		return
	}
	out := ent.envelopes[:0]
	for _, d := range ent.envelopes {
		if d == nil || d.CorrelationID != correlationID {
			out = append(out, d)
		}
	}
	ent.envelopes = out
}

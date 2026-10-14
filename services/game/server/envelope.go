package server

import (
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/route"
)

const (
	envelopeTimeout    = 5 * time.Second
	envelopeRetryDelay = 50 * time.Millisecond
)

func (gs *GameServer) Resend(d *route.Envelope) {
	if gs == nil || d == nil {
		return
	}
	time.AfterFunc(envelopeRetryDelay, func() {
		gs.retryEnvelope(d)
	})
}

func (gs *GameServer) CompleteEnvelope(correlationID uint64) {
	if gs == nil || correlationID == 0 {
		return
	}
	gs.envelopeMu.Lock()
	delete(gs.pendingEnvelopes, correlationID)
	gs.envelopeMu.Unlock()
}

func (gs *GameServer) SendToCharacter(caller *actor.PID, characterID uint32, inner interface{}) {
	if gs == nil || characterID == 0 || inner == nil {
		return
	}
	if gs.characterRuntime == nil {
		return
	}
	corr := gs.nextEnvelope.Add(1)
	if !gs.characterRuntime.Exists(characterID) {
		gs.tellEnvelopeResult(caller, corr, false, "not_registered")
		return
	}

	d := &route.Envelope{
		CorrelationID:    corr,
		CharacterID:      characterID,
		DeadlineUnixNano: clock.Now().Add(envelopeTimeout).UnixNano(),
		Caller:           caller,
		Inner:            inner,
	}

	gs.registerEnvelope(d)

	root := gs.GetRootContext()
	if root == nil {
		gs.failEnvelope(corr, "no_root")
		return
	}

	mapPID, ok := gs.characterRuntime.GetMapPID(characterID)
	if !ok {
		gs.failEnvelope(corr, "not_registered")
		return
	}
	if mapPID != nil {
		root.Send(mapPID, d)
		return
	}
	if err := gs.characterRuntime.queueEnvelope(characterID, d); err != nil {
		gs.failEnvelope(corr, "enqueue_failed")
	}
}

func (gs *GameServer) registerEnvelope(d *route.Envelope) {
	if d == nil {
		return
	}
	gs.envelopeMu.Lock()
	if gs.pendingEnvelopes == nil {
		gs.pendingEnvelopes = make(map[uint64]*route.Envelope)
	}
	gs.pendingEnvelopes[d.CorrelationID] = d
	gs.envelopeMu.Unlock()

	corr := d.CorrelationID
	time.AfterFunc(envelopeTimeout, func() {
		gs.failEnvelope(corr, "timeout")
	})
}

func (gs *GameServer) failEnvelope(correlation uint64, reason string) {
	if gs == nil || correlation == 0 {
		return
	}
	gs.envelopeMu.Lock()
	d, ok := gs.pendingEnvelopes[correlation]
	if ok {
		delete(gs.pendingEnvelopes, correlation)
	}
	gs.envelopeMu.Unlock()
	if !ok {
		return
	}
	if gs.characterRuntime != nil {
		gs.characterRuntime.removeEnvelope(d.CharacterID, correlation)
	}
	gs.tellEnvelopeResult(d.Caller, correlation, false, reason)
}

func (gs *GameServer) retryEnvelope(d *route.Envelope) {
	if gs == nil || d == nil || gs.characterRuntime == nil {
		return
	}
	if clock.Now().UnixNano() > d.DeadlineUnixNano {
		gs.failEnvelope(d.CorrelationID, "timeout")
		return
	}
	if !gs.characterRuntime.Exists(d.CharacterID) {
		gs.failEnvelope(d.CorrelationID, "gone")
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		gs.failEnvelope(d.CorrelationID, "no_root")
		return
	}
	mapPID, ok := gs.characterRuntime.GetMapPID(d.CharacterID)
	if !ok {
		gs.failEnvelope(d.CorrelationID, "gone")
		return
	}
	if mapPID == nil {
		time.AfterFunc(envelopeRetryDelay, func() {
			gs.retryEnvelope(d)
		})
		return
	}
	root.Send(mapPID, d)
}

func (gs *GameServer) flushEnvelopes(characterID uint32, mapPID *actor.PID) {
	if gs == nil || characterID == 0 || mapPID == nil || gs.characterRuntime == nil {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	for _, d := range gs.characterRuntime.takeEnvelopes(characterID) {
		if d == nil {
			continue
		}
		if clock.Now().UnixNano() > d.DeadlineUnixNano {
			gs.failEnvelope(d.CorrelationID, "timeout")
			continue
		}
		root.Send(mapPID, d)
	}
}

func (gs *GameServer) abandonEnvelopes(characterID uint32) {
	if gs == nil || characterID == 0 {
		return
	}
	var buffered []*route.Envelope
	if gs.characterRuntime != nil {
		buffered = gs.characterRuntime.takeEnvelopes(characterID)
	}
	for _, d := range buffered {
		if d != nil {
			gs.failEnvelope(d.CorrelationID, "gone")
		}
	}
	gs.envelopeMu.Lock()
	var rest []*route.Envelope
	for corr, d := range gs.pendingEnvelopes {
		if d != nil && d.CharacterID == characterID {
			rest = append(rest, d)
			delete(gs.pendingEnvelopes, corr)
		}
	}
	gs.envelopeMu.Unlock()
	for _, d := range rest {
		gs.tellEnvelopeResult(d.Caller, d.CorrelationID, false, "gone")
	}
}

func (gs *GameServer) tellEnvelopeResult(caller *actor.PID, correlation uint64, ok bool, reason string) {
	if caller == nil || correlation == 0 {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	root.Send(caller, &route.Result{
		CorrelationID: correlation,
		OK:            ok,
		Reason:        reason,
	})
}

package server

import (
	"github.com/boyism80/fm/core/clock"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/delivery"
)

const (
	deliveryTimeout    = 5 * time.Second
	deliveryRetryDelay = 50 * time.Millisecond
)

func (gs *GameServer) Redeliver(d *delivery.Envelope) {
	if gs == nil || d == nil {
		return
	}
	time.AfterFunc(deliveryRetryDelay, func() {
		gs.retryDelivery(d)
	})
}

func (gs *GameServer) CompleteDelivery(correlationID uint64) {
	if gs == nil || correlationID == 0 {
		return
	}
	gs.deliveryMu.Lock()
	delete(gs.pendingDeliveries, correlationID)
	gs.deliveryMu.Unlock()
}

func (gs *GameServer) Deliver(caller *actor.PID, characterID uint32, inner interface{}) {
	if gs == nil || characterID == 0 || inner == nil {
		return
	}
	if gs.characterRuntime == nil {
		return
	}
	corr := gs.nextDelivery.Add(1)
	if !gs.characterRuntime.Exists(characterID) {
		gs.tellDeliveryResult(caller, corr, false, "not_registered")
		return
	}

	d := &delivery.Envelope{
		CorrelationID:    corr,
		CharacterID:      characterID,
		DeadlineUnixNano: clock.Now().Add(deliveryTimeout).UnixNano(),
		Caller:           caller,
		Inner:            inner,
	}

	gs.registerDelivery(d)

	root := gs.GetRootContext()
	if root == nil {
		gs.failDelivery(corr, "no_root")
		return
	}

	mapPID, ok := gs.characterRuntime.GetMapPID(characterID)
	if !ok {
		gs.failDelivery(corr, "not_registered")
		return
	}
	if mapPID != nil {
		root.Send(mapPID, d)
		return
	}
	if err := gs.characterRuntime.queueDelivery(characterID, d); err != nil {
		gs.failDelivery(corr, "enqueue_failed")
	}
}

func (gs *GameServer) registerDelivery(d *delivery.Envelope) {
	if d == nil {
		return
	}
	gs.deliveryMu.Lock()
	if gs.pendingDeliveries == nil {
		gs.pendingDeliveries = make(map[uint64]*delivery.Envelope)
	}
	gs.pendingDeliveries[d.CorrelationID] = d
	gs.deliveryMu.Unlock()

	corr := d.CorrelationID
	time.AfterFunc(deliveryTimeout, func() {
		gs.failDelivery(corr, "timeout")
	})
}

func (gs *GameServer) failDelivery(correlation uint64, reason string) {
	if gs == nil || correlation == 0 {
		return
	}
	gs.deliveryMu.Lock()
	d, ok := gs.pendingDeliveries[correlation]
	if ok {
		delete(gs.pendingDeliveries, correlation)
	}
	gs.deliveryMu.Unlock()
	if !ok {
		return
	}
	if gs.characterRuntime != nil {
		gs.characterRuntime.removeDelivery(d.CharacterID, correlation)
	}
	gs.tellDeliveryResult(d.Caller, correlation, false, reason)
}

func (gs *GameServer) retryDelivery(d *delivery.Envelope) {
	if gs == nil || d == nil || gs.characterRuntime == nil {
		return
	}
	if clock.Now().UnixNano() > d.DeadlineUnixNano {
		gs.failDelivery(d.CorrelationID, "timeout")
		return
	}
	if !gs.characterRuntime.Exists(d.CharacterID) {
		gs.failDelivery(d.CorrelationID, "gone")
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		gs.failDelivery(d.CorrelationID, "no_root")
		return
	}
	mapPID, ok := gs.characterRuntime.GetMapPID(d.CharacterID)
	if !ok {
		gs.failDelivery(d.CorrelationID, "gone")
		return
	}
	if mapPID == nil {
		time.AfterFunc(deliveryRetryDelay, func() {
			gs.retryDelivery(d)
		})
		return
	}
	root.Send(mapPID, d)
}

func (gs *GameServer) flushDeliveries(characterID uint32, mapPID *actor.PID) {
	if gs == nil || characterID == 0 || mapPID == nil || gs.characterRuntime == nil {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	for _, d := range gs.characterRuntime.takeDeliveries(characterID) {
		if d == nil {
			continue
		}
		if clock.Now().UnixNano() > d.DeadlineUnixNano {
			gs.failDelivery(d.CorrelationID, "timeout")
			continue
		}
		root.Send(mapPID, d)
	}
}

func (gs *GameServer) abandonDeliveries(characterID uint32) {
	if gs == nil || characterID == 0 {
		return
	}
	var buffered []*delivery.Envelope
	if gs.characterRuntime != nil {
		buffered = gs.characterRuntime.takeDeliveries(characterID)
	}
	for _, d := range buffered {
		if d != nil {
			gs.failDelivery(d.CorrelationID, "gone")
		}
	}
	gs.deliveryMu.Lock()
	var rest []*delivery.Envelope
	for corr, d := range gs.pendingDeliveries {
		if d != nil && d.CharacterID == characterID {
			rest = append(rest, d)
			delete(gs.pendingDeliveries, corr)
		}
	}
	gs.deliveryMu.Unlock()
	for _, d := range rest {
		gs.tellDeliveryResult(d.Caller, d.CorrelationID, false, "gone")
	}
}

func (gs *GameServer) tellDeliveryResult(caller *actor.PID, correlation uint64, ok bool, reason string) {
	if caller == nil || correlation == 0 {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	root.Send(caller, &delivery.Result{
		CorrelationID: correlation,
		OK:            ok,
		Reason:        reason,
	})
}

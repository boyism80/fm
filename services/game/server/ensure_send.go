package server

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

const (
	ensureWallTimeout = 5 * time.Second
	ensureRetryDelay  = 50 * time.Millisecond
)

func (gs *GameServer) EnsureRedispatch(d *g_actor.EnsureDeliver) {
	if gs == nil || d == nil {
		return
	}
	time.AfterFunc(ensureRetryDelay, func() {
		gs.ensureRetryAttempt(d)
	})
}

func (gs *GameServer) EnsureComplete(correlationID uint64) {
	if gs == nil || correlationID == 0 {
		return
	}
	gs.ensureMu.Lock()
	delete(gs.ensurePending, correlationID)
	gs.ensureMu.Unlock()
}

func (gs *GameServer) EnsureSend(caller *actor.PID, characterID uint32, inner interface{}) {
	if gs == nil || characterID == 0 || inner == nil {
		return
	}
	if gs.characterRuntime == nil {
		return
	}
	corr := gs.ensureNext.Add(1)
	if !gs.characterRuntime.Exists(characterID) {
		gs.tellEnsureResult(caller, corr, false, "not_registered")
		return
	}

	d := &g_actor.EnsureDeliver{
		CorrelationID:    corr,
		CharacterID:      characterID,
		DeadlineUnixNano: time.Now().Add(ensureWallTimeout).UnixNano(),
		Caller:           caller,
		Inner:            inner,
	}

	gs.ensureRegisterPending(d)

	root := gs.GetRootContext()
	if root == nil {
		gs.ensureFailIfPending(corr, "no_root")
		return
	}

	mapPID, ok := gs.characterRuntime.GetMapPID(characterID)
	if !ok {
		gs.ensureFailIfPending(corr, "not_registered")
		return
	}
	if mapPID != nil {
		root.Send(mapPID, d)
		return
	}
	if err := gs.characterRuntime.enqueueEnsure(characterID, d); err != nil {
		gs.ensureFailIfPending(corr, "enqueue_failed")
	}
}

func (gs *GameServer) ensureRegisterPending(d *g_actor.EnsureDeliver) {
	if d == nil {
		return
	}
	gs.ensureMu.Lock()
	if gs.ensurePending == nil {
		gs.ensurePending = make(map[uint64]*g_actor.EnsureDeliver)
	}
	gs.ensurePending[d.CorrelationID] = d
	gs.ensureMu.Unlock()

	corr := d.CorrelationID
	time.AfterFunc(ensureWallTimeout, func() {
		gs.ensureFailIfPending(corr, "timeout")
	})
}

func (gs *GameServer) ensureFailIfPending(correlation uint64, reason string) {
	if gs == nil || correlation == 0 {
		return
	}
	gs.ensureMu.Lock()
	d, ok := gs.ensurePending[correlation]
	if ok {
		delete(gs.ensurePending, correlation)
	}
	gs.ensureMu.Unlock()
	if !ok {
		return
	}
	if gs.characterRuntime != nil {
		gs.characterRuntime.removeEnsureFromBuffer(d.CharacterID, correlation)
	}
	gs.tellEnsureResult(d.Caller, correlation, false, reason)
}

func (gs *GameServer) ensureRetryAttempt(d *g_actor.EnsureDeliver) {
	if gs == nil || d == nil || gs.characterRuntime == nil {
		return
	}
	if time.Now().UnixNano() > d.DeadlineUnixNano {
		gs.ensureFailIfPending(d.CorrelationID, "timeout")
		return
	}
	if !gs.characterRuntime.Exists(d.CharacterID) {
		gs.ensureFailIfPending(d.CorrelationID, "gone")
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		gs.ensureFailIfPending(d.CorrelationID, "no_root")
		return
	}
	mapPID, ok := gs.characterRuntime.GetMapPID(d.CharacterID)
	if !ok {
		gs.ensureFailIfPending(d.CorrelationID, "gone")
		return
	}
	if mapPID == nil {
		time.AfterFunc(ensureRetryDelay, func() {
			gs.ensureRetryAttempt(d)
		})
		return
	}
	root.Send(mapPID, d)
}

func (gs *GameServer) flushEnsureBuffer(characterID uint32, mapPID *actor.PID) {
	if gs == nil || characterID == 0 || mapPID == nil || gs.characterRuntime == nil {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	for _, d := range gs.characterRuntime.takeEnsureBuffer(characterID) {
		if d == nil {
			continue
		}
		if time.Now().UnixNano() > d.DeadlineUnixNano {
			gs.ensureFailIfPending(d.CorrelationID, "timeout")
			continue
		}
		root.Send(mapPID, d)
	}
}

func (gs *GameServer) ensureAbandonCharacter(characterID uint32) {
	if gs == nil || characterID == 0 {
		return
	}
	var buffered []*g_actor.EnsureDeliver
	if gs.characterRuntime != nil {
		buffered = gs.characterRuntime.takeEnsureBuffer(characterID)
	}
	for _, d := range buffered {
		if d != nil {
			gs.ensureFailIfPending(d.CorrelationID, "gone")
		}
	}
	gs.ensureMu.Lock()
	var rest []*g_actor.EnsureDeliver
	for corr, d := range gs.ensurePending {
		if d != nil && d.CharacterID == characterID {
			rest = append(rest, d)
			delete(gs.ensurePending, corr)
		}
	}
	gs.ensureMu.Unlock()
	for _, d := range rest {
		gs.tellEnsureResult(d.Caller, d.CorrelationID, false, "gone")
	}
}

func (gs *GameServer) tellEnsureResult(caller *actor.PID, correlation uint64, ok bool, reason string) {
	if caller == nil || correlation == 0 {
		return
	}
	root := gs.GetRootContext()
	if root == nil {
		return
	}
	root.Send(caller, &g_actor.EnsureResult{
		CorrelationID: correlation,
		OK:            ok,
		Reason:        reason,
	})
}

package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"google.golang.org/protobuf/proto"
)

type PartyEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	PartyID    uint32 `json:"party_id"`
	Revision   uint64 `json:"revision"`
	OccurredAt string `json:"occurred_at"`
}

type PartyContainer struct {
	gs             *GameServer
	worldID        uint32
	internalClient internal.InternalClient
	mu             sync.Mutex
	revisions      map[uint32]uint64
	snapshots      map[uint32]*internal.PartySnapshot
}

func NewPartyContainer(gs *GameServer, worldID uint32, ic internal.InternalClient) *PartyContainer {
	return &PartyContainer{
		gs:             gs,
		worldID:        worldID,
		internalClient: ic,
		revisions:      make(map[uint32]uint64),
		snapshots:      make(map[uint32]*internal.PartySnapshot),
	}
}

func (pc *PartyContainer) apply(evt PartyEventEnvelope) error {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	gs := pc.gs
	if evt.EventType == "log_onoff" {
		if err := pc.rehydrateLocked(evt.PartyID); err != nil {
			return err
		}
		if snapshot := pc.snapshots[evt.PartyID]; snapshot != nil && gs != nil {
			gs.SyncPartySnapshot(snapshot)
		}
		log.Printf("party consumer: applied type=log_onoff party_id=%d revision=%d", evt.PartyID, pc.revisions[evt.PartyID])
		return nil
	}
	switch evt.EventType {
	case "disbanded":
		var memberIDs []uint32
		if snapshot, ok := pc.snapshots[evt.PartyID]; ok && snapshot != nil {
			for _, m := range snapshot.GetMembers() {
				memberIDs = append(memberIDs, m.GetCharacterId())
			}
		}
		if len(memberIDs) == 0 {
			if err := pc.rehydrateLocked(evt.PartyID); err == nil {
				if snapshot, ok := pc.snapshots[evt.PartyID]; ok && snapshot != nil {
					for _, m := range snapshot.GetMembers() {
						memberIDs = append(memberIDs, m.GetCharacterId())
					}
				}
			}
		}
		delete(pc.revisions, evt.PartyID)
		delete(pc.snapshots, evt.PartyID)
		if gs != nil {
			gs.ClearPartyMembers(memberIDs)
		}
		log.Printf("party consumer: disbanded applied party_id=%d revision=%d", evt.PartyID, evt.Revision)
	default:
		if err := pc.rehydrateLocked(evt.PartyID); err != nil {
			return err
		}
		if snapshot := pc.snapshots[evt.PartyID]; snapshot != nil && gs != nil {
			gs.SyncPartySnapshot(snapshot)
		}
		log.Printf("party consumer: applied type=%s party_id=%d revision=%d", evt.EventType, evt.PartyID, pc.revisions[evt.PartyID])
	}
	return nil
}

func (pc *PartyContainer) rehydrateLocked(partyID uint32) error {
	if pc.internalClient == nil {
		return errors.New("internal client unavailable for rehydrate")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	reply, err := pc.internalClient.GetParty(ctx, &internal.GetPartyRequest{
		WorldId: pc.worldID,
		PartyId: partyID,
	})
	if err != nil {
		return fmt.Errorf("GetParty rehydrate failed: %w", err)
	}
	if !reply.GetFound() || reply.GetParty() == nil {
		delete(pc.revisions, partyID)
		delete(pc.snapshots, partyID)
		return nil
	}
	pc.revisions[partyID] = reply.GetParty().GetRevision()
	pc.snapshots[partyID] = reply.GetParty()
	log.Printf("party consumer: rehydrated party_id=%d revision=%d", partyID, pc.revisions[partyID])
	return nil
}

func (pc *PartyContainer) applyEmbeddedPartySnapshot(evt PartyEventEnvelope, snap *internal.PartySnapshot, notifySilent bool) (bool, error) {
	if pc == nil || snap == nil {
		return false, nil
	}
	partyID := evt.PartyID
	pc.mu.Lock()
	cloned, ok := proto.Clone(snap).(*internal.PartySnapshot)
	if !ok || cloned == nil {
		pc.mu.Unlock()
		return false, errors.New("party_snapshot: clone failed")
	}
	rev := cloned.GetRevision()
	pc.revisions[partyID] = rev
	pc.snapshots[partyID] = cloned
	pc.mu.Unlock()

	gs := pc.gs
	if gs != nil {
		gs.SyncPartySnapshot(cloned)
	}
	if notifySilent && gs != nil {
		gs.DeliverPartySilentFromSnapshot(cloned)
	}
	log.Printf("party consumer: applied embedded party snapshot party_id=%d revision=%d silent=%v", partyID, rev, notifySilent)
	return true, nil
}

func (pc *PartyContainer) ApplySnapshotFromLogin(snapshot *internal.PartySnapshot) {
	if pc == nil || snapshot == nil {
		return
	}
	partyID := snapshot.GetPartyId()
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.revisions[partyID] = snapshot.GetRevision()
	pc.snapshots[partyID] = snapshot
	if pc.gs != nil {
		pc.gs.SyncPartySnapshot(snapshot)
	}
	log.Printf("party consumer: hydrated from login party_id=%d revision=%d", partyID, snapshot.GetRevision())
}

func (pc *PartyContainer) CachedSnapshot(partyID uint32) *internal.PartySnapshot {
	if pc == nil {
		return nil
	}
	pc.mu.Lock()
	s := pc.snapshots[partyID]
	pc.mu.Unlock()
	if s == nil {
		return nil
	}
	cloned, ok := proto.Clone(s).(*internal.PartySnapshot)
	if !ok || cloned == nil {
		return nil
	}
	return cloned
}

package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
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
	parties        map[uint32]*entity.Party
}

func NewPartyContainer(gs *GameServer, worldID uint32, ic internal.InternalClient) *PartyContainer {
	return &PartyContainer{
		gs:             gs,
		worldID:        worldID,
		internalClient: ic,
		revisions:      make(map[uint32]uint64),
		parties:        make(map[uint32]*entity.Party),
	}
}

func (pc *PartyContainer) UpdateAsync(ctx actor.Context, evt PartyEventEnvelope) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if pc == nil {
		return p
	}
	partyID := evt.PartyID
	p.OnError(func(err error) {
	})

	if evt.EventType == "disbanded" {
		pc.mu.Lock()
		var memberIDs []uint32
		if party, ok := pc.parties[evt.PartyID]; ok && party != nil {
			memberIDs = entity.PartyMemberCharacterIDs(party.GetMembers())
		}
		if len(memberIDs) > 0 {
			delete(pc.revisions, evt.PartyID)
			delete(pc.parties, evt.PartyID)
			pc.ClearPartyMembers(memberIDs)
			pc.mu.Unlock()
			return p
		}
		pc.mu.Unlock()
	}

	if pc.internalClient == nil {
		p.Then(func(interface{}) (interface{}, error) {
			return nil, errors.New("party apply: internal client unavailable")
		})
		return p
	}

	async.ThenRPC(p,
		func(c context.Context) (*internal.GetPartyReply, error) {
			return pc.internalClient.GetParty(c, &internal.GetPartyRequest{
				WorldId: pc.worldID,
				PartyId: partyID,
			})
		},
		func(reply *internal.GetPartyReply) error {
			pc.mu.Lock()
			defer pc.mu.Unlock()
			pc.mergeGetPartyReplyLocked(partyID, reply)
			return pc.applyStateAfterRehydrateLocked(evt)
		},
	)
	return p
}

func (pc *PartyContainer) mergeGetPartyReplyLocked(partyID uint32, reply *internal.GetPartyReply) {
	if reply == nil || !reply.GetFound() || reply.GetParty() == nil {
		delete(pc.revisions, partyID)
		delete(pc.parties, partyID)
		return
	}
	ent := entity.PartyFromProto(pc.gs, reply.GetParty())
	if ent == nil {
		delete(pc.revisions, partyID)
		delete(pc.parties, partyID)
		return
	}
	stored := ent.Clone()
	if stored == nil {
		delete(pc.revisions, partyID)
		delete(pc.parties, partyID)
		return
	}
	pc.revisions[partyID] = stored.Revision
	pc.parties[partyID] = stored
}

func (pc *PartyContainer) applyStateAfterRehydrateLocked(evt PartyEventEnvelope) error {
	partyID := evt.PartyID
	if evt.EventType == "log_onoff" {
		if party := pc.parties[partyID]; party != nil {
			pc.Sync(party)
		}
		return nil
	}
	if evt.EventType == "disbanded" {
		var memberIDs []uint32
		if party, ok := pc.parties[partyID]; ok && party != nil {
			memberIDs = entity.PartyMemberCharacterIDs(party.GetMembers())
		}
		delete(pc.revisions, partyID)
		delete(pc.parties, partyID)
		pc.ClearPartyMembers(memberIDs)
		return nil
	}
	if party := pc.parties[partyID]; party != nil {
		pc.Sync(party)
	}
	return nil
}

func (pc *PartyContainer) applyEmbeddedParty(evt PartyEventEnvelope, partyPb *internal.Party, notifySilent bool) (bool, error) {
	if pc == nil || partyPb == nil {
		return false, nil
	}
	partyID := evt.PartyID
	ent := entity.PartyFromProto(pc.gs, partyPb)
	if ent == nil {
		return false, nil
	}
	cloned := ent.Clone()
	if cloned == nil {
		return false, errors.New("party: clone failed")
	}
	rev := cloned.Revision
	pc.mu.Lock()
	pc.revisions[partyID] = rev
	pc.parties[partyID] = cloned
	pc.mu.Unlock()

	pc.Sync(cloned)
	if notifySilent {
		pc.DeliverPartySilent(cloned)
	}
	return true, nil
}

func (pc *PartyContainer) Update(partyPb *internal.Party) {
	if pc == nil || partyPb == nil {
		return
	}
	ent := entity.PartyFromProto(pc.gs, partyPb)
	if ent == nil {
		return
	}
	stored := ent.Clone()
	if stored == nil {
		return
	}
	partyID := stored.PartyID
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.revisions[partyID] = stored.Revision
	pc.parties[partyID] = stored
	pc.Sync(stored)
}

func (pc *PartyContainer) Get(partyID uint32) *entity.Party {
	if pc == nil {
		return nil
	}
	pc.mu.Lock()
	s := pc.parties[partyID]
	pc.mu.Unlock()
	if s == nil {
		return nil
	}
	return s.Clone()
}

func (pc *PartyContainer) Sync(party *entity.Party) {
	if pc == nil || party == nil || pc.gs == nil {
		return
	}
	partyID := party.GetPartyId()
	for _, member := range party.GetMembers() {
		if member == nil {
			continue
		}
		cid := member.GetCharacterId()
		pid := partyID
		pc.gs.EnsureSend(nil, cid, &g_actor.SyncCharacterPartyState{
			CharacterID: cid,
			PartyID:     &pid,
		})
	}
}

func (pc *PartyContainer) ClearCharacterPartyID(characterID uint32) {
	if pc == nil || characterID == 0 || pc.gs == nil {
		return
	}
	pc.gs.EnsureSend(nil, characterID, &g_actor.SyncCharacterPartyState{
		CharacterID: characterID,
		PartyID:     nil,
	})
}

func (pc *PartyContainer) ClearPartyMembers(memberIDs []uint32) {
	if pc == nil || len(memberIDs) == 0 || pc.gs == nil {
		return
	}
	for _, cid := range memberIDs {
		pc.ClearCharacterPartyID(cid)
	}
}

func (pc *PartyContainer) DeliverPartySilent(party *entity.Party) {
	if pc == nil || pc.gs == nil || party == nil {
		return
	}
	gs := pc.gs
	respMembers := entity.PartyMembersToResponse(party.GetMembers())
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateSilent{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(gs.config.ChannelId),
			PartyID:     party.GetPartyId(),
			LeaderID:    party.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

// SendPartySilentAsync builds a Promise that sends party silent UI state to the character. Uses cache
// when warm; otherwise schedules GetParty via ThenRPC (from map actor Receive).
func (pc *PartyContainer) SendPartySilentAsync(ctx actor.Context, ch *entity.Character) *async.Promise {
	p := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if pc == nil || ch == nil || pc.gs == nil || ctx == nil {
		return p
	}
	partyIDPtr := ch.GetPartyID()
	if partyIDPtr == nil {
		return p
	}
	partyID := *partyIDPtr
	p.OnError(func(err error) {
		log.Printf("SendPartySilentAsync: GetParty char %d party %d: %v", ch.GetID(), partyID, err)
	})
	if ent := pc.Get(partyID); ent != nil {
		p.Then(func(interface{}) (interface{}, error) {
			pc.sendPartySilentToCharacter(ch, ent)
			return nil, nil
		})
		return p
	}
	if pc.internalClient == nil {
		return p
	}
	async.ThenRPC(p,
		func(c context.Context) (*internal.GetPartyReply, error) {
			return pc.internalClient.GetParty(c, &internal.GetPartyRequest{
				WorldId: pc.worldID,
				PartyId: partyID,
			})
		},
		func(reply *internal.GetPartyReply) error {
			if reply == nil || !reply.GetFound() || reply.GetParty() == nil {
				return nil
			}
			ent := entity.PartyFromProto(pc.gs, reply.GetParty())
			pc.sendPartySilentToCharacter(ch, ent)
			return nil
		},
	)
	return p
}

func (pc *PartyContainer) sendPartySilentToCharacter(ch *entity.Character, party *entity.Party) {
	if pc == nil || ch == nil || party == nil || pc.gs == nil {
		return
	}
	gs := pc.gs
	members := entity.PartyMembersToResponse(party.GetMembers())
	gs.EnsureSend(nil, ch.GetID(), &g_actor.DeliverPartyUpdateSilent{
		CharacterID: ch.GetID(),
		ForChannel:  int32(gs.config.ChannelId),
		PartyID:     party.GetPartyId(),
		LeaderID:    party.GetLeaderCharacterId(),
		Members:     members,
	})
}

func (pc *PartyContainer) PartyMemberIndex(characterID uint32, partyID *uint32) int {
	if partyID == nil || pc == nil {
		return 0
	}
	partyEnt := pc.Get(*partyID)
	if partyEnt == nil {
		return 0
	}
	for i, mem := range partyEnt.GetMembers() {
		if mem != nil && mem.GetCharacterId() == characterID {
			return i
		}
	}
	return 0
}

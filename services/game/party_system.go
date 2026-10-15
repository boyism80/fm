package game

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type PartyEventEnvelope struct {
	EventID    string `json:"event_id"`
	EventType  string `json:"event_type"`
	WorldID    uint32 `json:"world_id"`
	PartyID    uint32 `json:"party_id"`
	OccurredAt string `json:"occurred_at"`
}

type PartySystem struct {
	gw            GameWorld
	worldID       uint32
	mu            sync.Mutex
	parties       map[uint32]*Party
	memberUpdates map[uint32]*internal.UpdatePartyMemberRequest
}

func NewPartySystem(gw GameWorld) *PartySystem {
	return &PartySystem{
		gw:            gw,
		worldID:       gw.GetWorldID(),
		parties:       make(map[uint32]*Party),
		memberUpdates: make(map[uint32]*internal.UpdatePartyMemberRequest),
	}
}

func (s *PartySystem) Start() error { return nil }
func (s *PartySystem) Stop()        {}

func (s *PartySystem) UpdateAsync(ctx actor.Context, evt PartyEventEnvelope) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s == nil {
		return p
	}
	partyID := evt.PartyID
	p.OnError(func(err error) {
	})

	if evt.EventType == "disbanded" {
		s.mu.Lock()
		var memberIDs []uint32
		if party, ok := s.parties[evt.PartyID]; ok && party != nil {
			memberIDs = party.MemberCharacterIDs()
		}
		if len(memberIDs) > 0 {
			delete(s.parties, evt.PartyID)
			s.ClearPartyMembers(memberIDs)
			s.mu.Unlock()
			return p
		}
		s.mu.Unlock()
	}

	p.ThenRPC(
		func(c context.Context) (*internal.GetPartyReply, error) {
			return s.gw.Systems().Get[*InternalSystem]().client.GetParty(c, &internal.GetPartyRequest{
				WorldId: s.worldID,
				PartyId: partyID,
			})
		},
		func(reply *internal.GetPartyReply) error {
			s.mu.Lock()
			defer s.mu.Unlock()
			s.mergeGetPartyReplyLocked(partyID, reply)
			return s.syncAfterFetch(evt)
		},
	)
	return p
}

func (s *PartySystem) mergeGetPartyReplyLocked(partyID uint32, reply *internal.GetPartyReply) {
	if reply == nil || !reply.GetFound() || reply.GetParty() == nil {
		delete(s.parties, partyID)
		return
	}
	ent := PartyFromProto(s.gw, reply.GetParty())
	if ent == nil {
		delete(s.parties, partyID)
		return
	}
	stored := ent.Clone()
	if stored == nil {
		delete(s.parties, partyID)
		return
	}
	if prev := s.parties[partyID]; prev != nil && stored.UpdatedAt.Before(prev.UpdatedAt) {
		return
	}
	s.parties[partyID] = stored
}

func (s *PartySystem) syncAfterFetch(evt PartyEventEnvelope) error {
	partyID := evt.PartyID
	if evt.EventType == "log_onoff" {
		if party := s.parties[partyID]; party != nil {
			s.Sync(party)
		}
		return nil
	}
	if evt.EventType == "disbanded" {
		var memberIDs []uint32
		if party, ok := s.parties[partyID]; ok && party != nil {
			memberIDs = party.MemberCharacterIDs()
		}
		delete(s.parties, partyID)
		s.ClearPartyMembers(memberIDs)
		return nil
	}
	if party := s.parties[partyID]; party != nil {
		s.Sync(party)
	}
	return nil
}

func (s *PartySystem) ApplyPartySnapshot(evt PartyEventEnvelope, partyPb *internal.Party, notifySilent bool) (bool, error) {
	if s == nil || partyPb == nil {
		return false, nil
	}
	partyID := evt.PartyID
	ent := PartyFromProto(s.gw, partyPb)
	if ent == nil {
		return false, nil
	}
	cloned := ent.Clone()
	if cloned == nil {
		return false, errors.New("party: clone failed")
	}
	s.mu.Lock()
	if prev := s.parties[partyID]; prev != nil && cloned.UpdatedAt.Before(prev.UpdatedAt) {
		s.mu.Unlock()
		return false, nil
	}
	s.parties[partyID] = cloned
	s.mu.Unlock()

	s.Sync(cloned)
	if notifySilent {
		s.BroadcastPartySilent(cloned)
	}
	return true, nil
}

func (s *PartySystem) Update(partyPb *internal.Party) {
	if s == nil || partyPb == nil {
		return
	}
	ent := PartyFromProto(s.gw, partyPb)
	if ent == nil {
		return
	}
	stored := ent.Clone()
	if stored == nil {
		return
	}
	partyID := stored.PartyID
	s.mu.Lock()
	defer s.mu.Unlock()
	if prev := s.parties[partyID]; prev != nil && stored.UpdatedAt.Before(prev.UpdatedAt) {
		return
	}
	s.parties[partyID] = stored
	s.Sync(stored)
}

func (s *PartySystem) Get(partyID uint32) *Party {
	if s == nil {
		return nil
	}
	s.mu.Lock()
	party := s.parties[partyID]
	s.mu.Unlock()
	if party == nil {
		return nil
	}
	return party.Clone()
}

func (s *PartySystem) Sync(party *Party) {
	if s == nil || party == nil || s.gw == nil {
		return
	}
	partyID := party.GetPartyId()
	for _, member := range party.GetMembers() {
		if member == nil {
			continue
		}
		cid := member.GetCharacterId()
		pid := partyID
		s.gw.Systems().Get[*DispatchSystem]().SendToCharacter(nil, cid, &SyncCharacterPartyState{
			CharacterID: cid,
			PartyID:     &pid,
		})
	}
}

func (s *PartySystem) ClearPartyMembers(memberIDs []uint32) {
	if s == nil || len(memberIDs) == 0 || s.gw == nil {
		return
	}
	for _, cid := range memberIDs {
		if cid == 0 {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().SendToCharacter(nil, cid, &SyncCharacterPartyState{
			CharacterID: cid,
			PartyID:     nil,
		})
	}
}

func (s *PartySystem) BroadcastPartySilent(party *Party) {
	if s == nil || s.gw == nil || party == nil {
		return
	}
	respMembers := party.BuildS2C()
	forChannel := s.gw.GetChannelID()
	partyID := party.GetPartyId()
	leaderID := party.GetLeaderCharacterId()
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(m.GetCharacterId(), func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnPartyUpdateSilent(ch, forChannel, partyID, leaderID, respMembers)
		})
	}
}

func (s *PartySystem) SendPartySilentAsync(ctx actor.Context, ch *Character) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s == nil || ch == nil || s.gw == nil || ctx == nil {
		return p
	}
	partyIDPtr := ch.Party.ID()
	if partyIDPtr == nil {
		return p
	}
	partyID := *partyIDPtr
	p.OnError(func(err error) {
		log.Printf("SendPartySilentAsync: GetParty char %d party %d: %v", ch.GetID(), partyID, err)
	})
	if ent := s.Get(partyID); ent != nil {
		p.Do(func() error {
			s.sendPartySilentToCharacter(ch, ent)
			return nil
		})
		return p
	}
	p.ThenRPC(
		func(c context.Context) (*internal.GetPartyReply, error) {
			return s.gw.Systems().Get[*InternalSystem]().client.GetParty(c, &internal.GetPartyRequest{
				WorldId: s.worldID,
				PartyId: partyID,
			})
		},
		func(reply *internal.GetPartyReply) error {
			if reply == nil || !reply.GetFound() || reply.GetParty() == nil {
				return nil
			}
			ent := PartyFromProto(s.gw, reply.GetParty())
			s.sendPartySilentToCharacter(ch, ent)
			return nil
		},
	)
	return p
}

func (s *PartySystem) sendPartySilentToCharacter(ch *Character, party *Party) {
	if s == nil || ch == nil || party == nil || s.gw == nil {
		return
	}
	members := party.BuildS2C()
	forChannel := s.gw.GetChannelID()
	partyID := party.GetPartyId()
	leaderID := party.GetLeaderCharacterId()
	s.gw.Systems().Get[*DispatchSystem]().Visit(ch.GetID(), func(_ actor.Context, ch *Character) {
		if ch == nil {
			return
		}
		ch.Listener.OnPartyUpdateSilent(ch, forChannel, partyID, leaderID, members)
	})
}

func (s *PartySystem) PartyMemberIndex(characterID uint32, partyID *uint32) int {
	if partyID == nil || s == nil {
		return 0
	}
	partyEnt := s.Get(*partyID)
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

package game

import (
	"context"
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
	PartyPB    string `json:"party_pb"`
	Revision   uint64 `json:"revision"`
}

type PartySystem struct {
	gw            GameWorld
	worldID       uint32
	mu            sync.Mutex
	parties       map[uint32]*Party
	memberUpdates map[uint32]*PartyMember
}

func NewPartySystem(gw GameWorld) *PartySystem {
	return &PartySystem{
		gw:            gw,
		worldID:       gw.GetWorldID(),
		parties:       make(map[uint32]*Party),
		memberUpdates: make(map[uint32]*PartyMember),
	}
}

func (s *PartySystem) Start() error { return nil }
func (s *PartySystem) Stop()        {}

func (s *PartySystem) Update(partyPb *internal.Party) bool {
	if partyPb == nil {
		return false
	}
	if s.Get(partyPb.GetPartyId()) == nil {
		found := false
		for _, m := range partyPb.GetMembers() {
			if s.gw.Systems().Get[*DispatchSystem]().Exists(m.GetCharacterId()) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	party := PartyFromProto(s.gw, partyPb)
	if party == nil {
		return false
	}
	s.mu.Lock()
	if prev := s.parties[party.partyID]; prev != nil && party.revision < prev.revision {
		s.mu.Unlock()
		return false
	}
	s.parties[party.partyID] = party
	s.mu.Unlock()

	s.Sync(party)
	return true
}

func (s *PartySystem) Disband(partyID uint32, revision uint64) {
	s.mu.Lock()
	var memberIDs []uint32
	if party := s.parties[partyID]; party != nil {
		memberIDs = party.MemberCharacterIDs()
	}
	delete(s.parties, partyID)
	s.mu.Unlock()

	for _, cid := range memberIDs {
		s.gw.Systems().Get[*DispatchSystem]().SendToCharacter(nil, cid, &SyncCharacterPartyState{
			CharacterID: cid,
			PartyID:     partyID,
			Revision:    revision,
			Clear:       true,
		})
	}
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
	for _, member := range party.GetMembers() {
		if member == nil {
			continue
		}
		cid := member.GetCharacterId()
		s.gw.Systems().Get[*DispatchSystem]().SendToCharacter(nil, cid, &SyncCharacterPartyState{
			CharacterID: cid,
			PartyID:     party.partyID,
			Revision:    party.revision,
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
		if m == nil {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(m.GetCharacterId(), func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnPartyUpdateSilent(ch, forChannel, partyID, leaderID, respMembers)
		})
	}
}

func (s *PartySystem) SendPartySilentAsync(ctx actor.Context, ch *Character) *async.Task {
	p := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s == nil || ch == nil || s.gw == nil || ctx == nil {
		return p
	}
	partyIDPtr := ch.party.ID()
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
		ch.listener.OnPartyUpdateSilent(ch, forChannel, partyID, leaderID, members)
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

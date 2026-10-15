package game

import (
	"github.com/asynkron/protoactor-go/actor"
	c_actor "github.com/boyism80/fm/core/actor"
	"github.com/boyism80/fm/protocol/response"
)

func (s *PartySystem) BroadcastMemberJoined(party *Party, joinedCharacterID uint32) {
	if s == nil || s.gw == nil || party == nil || joinedCharacterID == 0 {
		return
	}
	members := party.GetMembers()
	if len(members) == 0 {
		return
	}
	joinName := ""
	for _, m := range members {
		if m != nil && m.GetCharacterId() == joinedCharacterID {
			joinName = m.GetCharacterName()
			break
		}
	}
	if joinName == "" {
		return
	}
	respMembers := party.BuildS2C()
	for _, m := range members {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().SendToCharacter(nil, m.GetCharacterId(), &PartyMemberJoined{
			CharacterID:     m.GetCharacterId(),
			ForChannel:      s.gw.GetChannelID(),
			PartyID:         party.GetPartyId(),
			JoinCharacterID: joinedCharacterID,
			JoinName:        joinName,
			LeaderID:        party.GetLeaderCharacterId(),
			Members:         respMembers,
		})
	}
}

func (s *PartySystem) BroadcastMemberLeft(prev, current *Party, targetCharacterID uint32, expelled bool) {
	if s == nil || s.gw == nil || prev == nil || targetCharacterID == 0 {
		return
	}
	targetName := ""
	oldMembers := prev.GetMembers()
	for _, m := range oldMembers {
		if m != nil && m.GetCharacterId() == targetCharacterID {
			targetName = m.GetCharacterName()
			break
		}
	}
	if targetName == "" {
		return
	}
	partyID := prev.GetPartyId()
	leaderID := prev.GetLeaderCharacterId()
	respMembers := []response.PartyMemberStatus{}
	if current != nil {
		partyID = current.GetPartyId()
		leaderID = current.GetLeaderCharacterId()
		respMembers = current.BuildS2C()
	}

	oldIDs := prev.MemberCharacterIDs()
	if root := s.gw.Systems().Get[*c_actor.ActorRegistry]().Root(); root != nil {
		seen := make(map[string]struct{})
		for _, cid := range oldIDs {
			mapPID, ok := s.gw.Systems().Get[*DispatchSystem]().GetMapPID(cid)
			if !ok || mapPID == nil {
				continue
			}
			key := mapPID.String()
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			root.Send(mapPID, &PartyMemberLeft{
				LeaverID: targetCharacterID,
			})
		}
	}

	forChannel := s.gw.GetChannelID()
	for _, m := range oldMembers {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(m.GetCharacterId(), func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			if expelled {
				ch.Listener.OnPartyUpdateExpel(ch, forChannel, partyID, targetCharacterID, targetName, leaderID, respMembers)
				return
			}
			ch.Listener.OnPartyUpdateLeave(ch, forChannel, partyID, targetCharacterID, targetName, leaderID, respMembers)
		})
	}
}

func (s *PartySystem) BroadcastDisbanded(prev *Party, leaderCharacterID uint32) {
	if s == nil || s.gw == nil || prev == nil || leaderCharacterID == 0 {
		return
	}

	formerIDs := prev.MemberCharacterIDs()
	if len(formerIDs) > 0 {
		if root := s.gw.Systems().Get[*c_actor.ActorRegistry]().Root(); root != nil {
			seen := make(map[string]struct{})
			for _, cid := range formerIDs {
				mapPID, ok := s.gw.Systems().Get[*DispatchSystem]().GetMapPID(cid)
				if !ok || mapPID == nil {
					continue
				}
				key := mapPID.String()
				if _, dup := seen[key]; dup {
					continue
				}
				seen[key] = struct{}{}
				root.Send(mapPID, &PartyDisband{
					FormerMemberIDs: formerIDs,
				})
			}
		}
	}

	partyID := prev.GetPartyId()
	for _, m := range prev.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(m.GetCharacterId(), func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnPartyUpdateDisband(ch, partyID, leaderCharacterID)
		})
	}
}

func (s *PartySystem) BroadcastLeaderChanged(party *Party, newLeaderCharacterID uint32, disconnected bool) {
	if s == nil || s.gw == nil || party == nil || newLeaderCharacterID == 0 {
		return
	}
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(m.GetCharacterId(), func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnPartyUpdateLeaderChange(ch, newLeaderCharacterID, disconnected)
		})
	}
}

func (s *PartySystem) BroadcastLogOnOff(party *Party) {
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
			ch.Listener.OnPartyUpdateLogOnOff(ch, forChannel, partyID, leaderID, respMembers)
		})
	}
}

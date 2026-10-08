package server

import (
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
)

func (pc *PartyCache) BroadcastMemberJoined(party *entity.Party, joinedCharacterID uint32) {
	if pc == nil || pc.gs == nil || party == nil || joinedCharacterID == 0 {
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
	respMembers := entity.PartyMembersToResponse(members)
	for _, m := range members {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		pc.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateJoin{
			CharacterID:     m.GetCharacterId(),
			ForChannel:      int32(pc.gs.config.ChannelId),
			PartyID:         party.GetPartyId(),
			JoinCharacterID: joinedCharacterID,
			JoinName:        joinName,
			LeaderID:        party.GetLeaderCharacterId(),
			Members:         respMembers,
		})
	}
}

func (pc *PartyCache) BroadcastMemberLeft(prev, current *entity.Party, targetCharacterID uint32, expelled bool) {
	if pc == nil || pc.gs == nil || prev == nil || targetCharacterID == 0 {
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
	var members []*entity.PartyMember
	if current != nil {
		partyID = current.GetPartyId()
		leaderID = current.GetLeaderCharacterId()
		members = current.GetMembers()
	}
	respMembers := entity.PartyMembersToResponse(members)

	if pc.gs.characterRuntime != nil {
		oldIDs := entity.PartyMemberCharacterIDs(oldMembers)
		if root := pc.gs.GetRootContext(); root != nil {
			seen := make(map[string]struct{})
			for _, cid := range oldIDs {
				mapPID, ok := pc.gs.characterRuntime.GetMapPID(cid)
				if !ok || mapPID == nil {
					continue
				}
				key := mapPID.String()
				if _, dup := seen[key]; dup {
					continue
				}
				seen[key] = struct{}{}
				root.Send(mapPID, &g_actor.PartyMemberLeft{
					LeaverID: targetCharacterID,
				})
			}
		}
	}

	for _, m := range oldMembers {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		pc.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeave{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(pc.gs.config.ChannelId),
			PartyID:     partyID,
			TargetID:    targetCharacterID,
			TargetName:  targetName,
			LeaderID:    leaderID,
			Members:     respMembers,
			Expelled:    expelled,
		})
	}
}

func (pc *PartyCache) BroadcastDisbanded(prev *entity.Party, leaderCharacterID uint32) {
	if pc == nil || pc.gs == nil || prev == nil || leaderCharacterID == 0 {
		return
	}

	if pc.gs.characterRuntime != nil {
		formerIDs := entity.PartyMemberCharacterIDs(prev.GetMembers())
		if len(formerIDs) > 0 {
			if root := pc.gs.GetRootContext(); root != nil {
				seen := make(map[string]struct{})
				for _, cid := range formerIDs {
					mapPID, ok := pc.gs.characterRuntime.GetMapPID(cid)
					if !ok || mapPID == nil {
						continue
					}
					key := mapPID.String()
					if _, dup := seen[key]; dup {
						continue
					}
					seen[key] = struct{}{}
					root.Send(mapPID, &g_actor.PartyDisband{
						FormerMemberIDs: formerIDs,
					})
				}
			}
		}
	}

	for _, m := range prev.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		pc.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateDisband{
			CharacterID: m.GetCharacterId(),
			PartyID:     prev.GetPartyId(),
			LeaderID:    leaderCharacterID,
		})
	}
}

func (pc *PartyCache) BroadcastLeaderChanged(party *entity.Party, newLeaderCharacterID uint32, disconnected bool) {
	if pc == nil || pc.gs == nil || party == nil || newLeaderCharacterID == 0 {
		return
	}
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		pc.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeaderChange{
			CharacterID:          m.GetCharacterId(),
			NewLeaderCharacterID: newLeaderCharacterID,
			ByDisconnect:         disconnected,
		})
	}
}

func (pc *PartyCache) BroadcastLogOnOff(party *entity.Party) {
	if pc == nil || pc.gs == nil || party == nil {
		return
	}
	respMembers := entity.PartyMembersToResponse(party.GetMembers())
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		pc.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLogOnOff{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(pc.gs.config.ChannelId),
			PartyID:     party.GetPartyId(),
			LeaderID:    party.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

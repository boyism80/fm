package server

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type partyMqHandler struct {
	gs *GameServer
}

func (h *partyMqHandler) sendInvite(targetCharacterID uint32, partyID uint32, inviterName string, partySearch bool) {
	if h.gs == nil || targetCharacterID == 0 {
		return
	}
	h.gs.EnsureSend(nil, targetCharacterID, &g_actor.DeliverPartyInvite{
		CharacterID: targetCharacterID,
		PartyID:     partyID,
		InviterName: inviterName,
		PartySearch: partySearch,
	})
}

func (h *partyMqHandler) sendDenyStatus(targetCharacterID uint32, action uint8, deniedCharacterName string) {
	if h.gs == nil || targetCharacterID == 0 {
		return
	}
	h.gs.EnsureSend(nil, targetCharacterID, &g_actor.DeliverPartyStatusMessage{
		CharacterID: targetCharacterID,
		Code:        pconst.PartyStatusCode(action),
		Name:        deniedCharacterName,
	})
}

func (h *partyMqHandler) notifyMemberLeftToMaps(leaverID uint32, prev *entity.Party) {
	if h.gs == nil || h.gs.characterRuntime == nil || leaverID == 0 || prev == nil {
		return
	}
	oldIDs := entity.PartyMemberCharacterIDs(prev.GetMembers())
	root := h.gs.GetRootContext()
	if root == nil {
		return
	}
	seen := make(map[string]struct{})
	for _, cid := range oldIDs {
		mapPID, ok := h.gs.characterRuntime.GetMapPID(cid)
		if !ok || mapPID == nil {
			continue
		}
		key := mapPID.String()
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		root.Send(mapPID, &g_actor.PartyMemberLeft{
			LeaverID: leaverID,
		})
	}
}

func (h *partyMqHandler) notifyDisbandToMaps(prev *entity.Party) {
	if h.gs == nil || h.gs.characterRuntime == nil || prev == nil {
		return
	}
	formerIDs := entity.PartyMemberCharacterIDs(prev.GetMembers())
	if len(formerIDs) == 0 {
		return
	}
	root := h.gs.GetRootContext()
	if root == nil {
		return
	}
	seen := make(map[string]struct{})
	for _, cid := range formerIDs {
		mapPID, ok := h.gs.characterRuntime.GetMapPID(cid)
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

func (h *partyMqHandler) sendJoinUpdate(party *entity.Party, joinedCharacterID uint32) {
	if h.gs == nil || party == nil || joinedCharacterID == 0 {
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
	respMembers := entity.PartyMembersToResponse(members)
	if joinName == "" {
		return
	}
	for _, m := range members {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		h.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateJoin{
			CharacterID:     m.GetCharacterId(),
			ForChannel:      int32(h.gs.config.ChannelId),
			PartyID:         party.GetPartyId(),
			JoinCharacterID: joinedCharacterID,
			JoinName:        joinName,
			LeaderID:        party.GetLeaderCharacterId(),
			Members:         respMembers,
		})
	}
}

func (h *partyMqHandler) sendLeaveUpdate(prev, current *entity.Party, targetCharacterID uint32, expelled bool) {
	if h.gs == nil || prev == nil || targetCharacterID == 0 {
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
	var (
		partyID  = prev.GetPartyId()
		leaderID = prev.GetLeaderCharacterId()
		members  []*entity.PartyMember
	)
	if current != nil {
		partyID = current.GetPartyId()
		leaderID = current.GetLeaderCharacterId()
		members = current.GetMembers()
	}
	respMembers := entity.PartyMembersToResponse(members)
	h.notifyMemberLeftToMaps(targetCharacterID, prev)
	for _, m := range oldMembers {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		h.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeave{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(h.gs.config.ChannelId),
			PartyID:     partyID,
			TargetID:    targetCharacterID,
			TargetName:  targetName,
			LeaderID:    leaderID,
			Members:     respMembers,
			Expelled:    expelled,
		})
	}
}

func (h *partyMqHandler) sendDisbandUpdate(prev *entity.Party, leaderCharacterID uint32) {
	if h.gs == nil || prev == nil || leaderCharacterID == 0 {
		return
	}
	h.notifyDisbandToMaps(prev)
	for _, m := range prev.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		h.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateDisband{
			CharacterID: m.GetCharacterId(),
			PartyID:     prev.GetPartyId(),
			LeaderID:    leaderCharacterID,
		})
	}
}

func (h *partyMqHandler) sendLeaderChange(party *entity.Party, newLeaderCharacterID uint32, byDisconnect bool) {
	if h.gs == nil || party == nil || newLeaderCharacterID == 0 {
		return
	}
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		h.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLeaderChange{
			CharacterID:          m.GetCharacterId(),
			NewLeaderCharacterID: newLeaderCharacterID,
			ByDisconnect:         byDisconnect,
		})
	}
}

func (h *partyMqHandler) sendLogOnOff(party *entity.Party) {
	if h.gs == nil || party == nil {
		return
	}
	respMembers := entity.PartyMembersToResponse(party.GetMembers())
	for _, m := range party.GetMembers() {
		if m == nil || m.GetCharacterId() == 0 {
			continue
		}
		h.gs.EnsureSend(nil, m.GetCharacterId(), &g_actor.DeliverPartyUpdateLogOnOff{
			CharacterID: m.GetCharacterId(),
			ForChannel:  int32(h.gs.config.ChannelId),
			PartyID:     party.GetPartyId(),
			LeaderID:    party.GetLeaderCharacterId(),
			Members:     respMembers,
		})
	}
}

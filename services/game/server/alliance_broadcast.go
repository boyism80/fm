package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

func (ac *AllianceContainer) BroadcastCreate(alliancePb *internal.Alliance) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := ac.guildInfos(alliancePb)
	membershipGuilds := entity.AllianceMembershipChangeGuildsFromProto(alliancePb)
	if len(membershipGuilds) == 0 {
		membershipGuilds = ac.membershipGuilds(alliancePb)
	}
	for _, memberID := range ac.recipientIDs(alliancePb) {
		ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceCreate{
			CharacterID:      memberID,
			Info:             info,
			Guilds:           guilds,
			MembershipGuilds: membershipGuilds,
		})
	}
}

func (ac *AllianceContainer) guildInfos(alliancePb *internal.Alliance) []*dto.GuildInfo {
	if guilds := entity.AllianceCreateGuildsFromProto(alliancePb); len(guilds) > 0 {
		return guilds
	}
	if ac.gs == nil || alliancePb == nil {
		return nil
	}
	allianceID := alliancePb.GetAllianceId()
	var out []*dto.GuildInfo
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := ac.gs.guild.Get(guildID)
		if ent == nil {
			continue
		}
		pb := ent.ToProto()
		if pb == nil {
			continue
		}
		if pb.AllianceId == nil {
			id := allianceID
			pb.AllianceId = &id
		}
		if info := entity.GuildInfoFromProto(pb); info != nil {
			out = append(out, info)
		}
	}
	return out
}

func (ac *AllianceContainer) membershipGuilds(alliancePb *internal.Alliance) []dto.AllianceMembershipChangeGuild {
	if ac.gs == nil || alliancePb == nil {
		return nil
	}
	var out []dto.AllianceMembershipChangeGuild
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := ac.gs.guild.Get(guildID)
		if ent == nil {
			continue
		}
		pb := ent.ToProto()
		if pb == nil {
			continue
		}
		block, ok := entity.AllianceMembershipChangeGuildFromProto(pb)
		if ok {
			out = append(out, block)
		}
	}
	return out
}

func (ac *AllianceContainer) recipientIDs(alliancePb *internal.Alliance) []uint32 {
	if alliancePb == nil {
		return nil
	}
	seen := make(map[uint32]struct{})
	var ids []uint32
	add := func(characterID uint32) {
		if characterID == 0 {
			return
		}
		if _, ok := seen[characterID]; ok {
			return
		}
		seen[characterID] = struct{}{}
		ids = append(ids, characterID)
	}
	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m != nil {
				add(m.GetCharacterId())
			}
		}
	}
	if len(ids) > 0 {
		return ids
	}
	if ac.gs == nil {
		return ids
	}
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := ac.gs.guild.Get(guildID)
		if ent == nil {
			continue
		}
		for _, m := range ent.GetMembers() {
			if m != nil {
				add(m.GetCharacterId())
			}
		}
	}
	return ids
}

func (ac *AllianceContainer) forEachOnlineAllianceMember(alliancePb *internal.Alliance, fn func(memberID uint32)) {
	if ac.gs == nil || alliancePb == nil || fn == nil {
		return
	}
	seen := make(map[uint32]struct{})
	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 {
				continue
			}
			if _, dup := seen[memberID]; dup {
				continue
			}
			seen[memberID] = struct{}{}
			if ac.gs.characterRuntime == nil || !ac.gs.characterRuntime.Exists(memberID) {
				continue
			}
			fn(memberID)
		}
	}
}

func (ac *AllianceContainer) deliverToOnlineMember(memberID uint32, msg interface{}) {
	if ac.gs == nil || memberID == 0 || msg == nil {
		return
	}
	if ac.gs.characterRuntime == nil || !ac.gs.characterRuntime.Exists(memberID) {
		return
	}
	ac.gs.EnsureSend(nil, memberID, msg)
}

func (ac *AllianceContainer) BroadcastNoticeChanged(alliancePb *internal.Alliance, notice string) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceNoticeChanged{
			CharacterID: memberID,
			Info:        info,
		})
	})
	if notice == "" {
		return
	}
	msg := gameconst.AllianceNoticeChangedMessagePrefix + notice
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.gs.EnsureSend(nil, memberID, &g_actor.DeliverMessage{
			CharacterID: memberID,
			MessageType: gameconst.MsgPopup,
			Message:     msg,
		})
	})
}

func (ac *AllianceContainer) BroadcastLeaderChanged(alliancePb *internal.Alliance, oldLeaderID uint32, newLeaderID uint32) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	allianceID := alliancePb.GetAllianceId()
	guilds := ac.guildInfos(alliancePb)
	leaderNotice := ""
	if name := ac.allianceCharacterName(alliancePb, newLeaderID); name != "" {
		leaderNotice = name + gameconst.AllianceLeaderChangedMessageSuffix
	}
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceLeaderChanged{
			CharacterID: memberID,
			AllianceID:  allianceID,
			OldLeaderID: oldLeaderID,
			NewLeaderID: newLeaderID,
			Info:        info,
			Guilds:      guilds,
		})
		if leaderNotice != "" {
			ac.gs.EnsureSend(nil, memberID, &g_actor.DeliverMessage{
				CharacterID: memberID,
				MessageType: gameconst.MsgPinkText,
				Message:     leaderNotice,
			})
		}
	})
}

func (ac *AllianceContainer) BroadcastInfoUpdate(alliancePb *internal.Alliance) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceInfoBroadcast{
			CharacterID: memberID,
			Info:        info,
		})
	})
}

func (ac *AllianceContainer) broadcastAllianceStateRefresh(alliancePb *internal.Alliance) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := ac.guildInfos(alliancePb)
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceMemberRankChanged{
			CharacterID: memberID,
			Info:        info,
			Guilds:      guilds,
		})
	})
}

func (ac *AllianceContainer) allianceCharacterName(alliancePb *internal.Alliance, characterID uint32) string {
	if alliancePb == nil || characterID == 0 {
		return ""
	}
	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m != nil && m.GetCharacterId() == characterID {
				return m.GetCharacterName()
			}
		}
	}
	return ""
}

func (ac *AllianceContainer) BroadcastMemberRankChanged(alliancePb *internal.Alliance, characterID uint32, allianceRank uint32) {
	if ac.gs == nil || alliancePb == nil || characterID == 0 {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	if allianceRank == 0 {
		for _, g := range alliancePb.GetGuilds() {
			if g == nil {
				continue
			}
			for _, m := range g.GetMembers() {
				if m != nil && m.GetCharacterId() == characterID && m.AllianceRank != nil {
					allianceRank = m.GetAllianceRank()
					break
				}
			}
			if allianceRank != 0 {
				break
			}
		}
	}
	if allianceRank == 0 {
		return
	}
	ac.broadcastAllianceStateRefresh(alliancePb)
}

func (ac *AllianceContainer) BroadcastGuildAdded(alliancePb *internal.Alliance, addedGuildPb *internal.Guild) {
	if ac.gs == nil || alliancePb == nil || addedGuildPb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	addedInfo := entity.GuildInfoFromProto(addedGuildPb)
	if info == nil || addedInfo == nil {
		return
	}
	guilds := ac.guildInfos(alliancePb)
	newGuildID := addedGuildPb.GetGuildId()
	members := entity.AllianceGuildMemberRanksFromProto(addedGuildPb)
	var membershipGuild dto.AllianceMembershipChangeGuild
	hasMembership := false
	if block, ok := entity.AllianceMembershipChangeGuildFromProto(addedGuildPb); ok {
		membershipGuild = block
		hasMembership = true
	}

	seen := make(map[uint32]struct{})
	deliver := func(memberID uint32, joining bool) {
		if memberID == 0 {
			return
		}
		if _, dup := seen[memberID]; dup {
			return
		}
		seen[memberID] = struct{}{}
		msg := &g_actor.DeliverAllianceGuildAdded{
			CharacterID:   memberID,
			Info:          info,
			Guilds:        guilds,
			NewGuildID:    newGuildID,
			AddedGuild:    addedInfo,
			Members:       members,
			Joining:       joining,
			HasMembership: hasMembership,
		}
		if hasMembership {
			msg.MembershipGuild = membershipGuild
		}
		ac.deliverToOnlineMember(memberID, msg)
	}

	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			deliver(m.GetCharacterId(), false)
		}
	}
	for _, m := range addedGuildPb.GetMembers() {
		if m == nil {
			continue
		}
		deliver(m.GetCharacterId(), true)
	}
	ac.broadcastAllianceStateRefresh(alliancePb)
}

func (ac *AllianceContainer) BroadcastGuildLeft(alliancePb *internal.Alliance, removedGuildPb *internal.Guild, expelled bool) {
	if ac.gs == nil {
		return
	}
	if alliancePb != nil {
		ac.Update(alliancePb)
	}
	if removedGuildPb != nil {
		ac.gs.guild.Update(removedGuildPb)
	}
	if alliancePb == nil || removedGuildPb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	removedInfo := entity.GuildInfoFromProto(removedGuildPb)
	if info == nil || removedInfo == nil {
		return
	}
	removedGuildID := removedGuildPb.GetGuildId()
	removedMembers := entity.AllianceGuildMemberRanksFromProto(removedGuildPb)

	seen := make(map[uint32]struct{})
	deliver := func(memberID uint32, leaving bool) {
		if memberID == 0 {
			return
		}
		if _, dup := seen[memberID]; dup {
			return
		}
		seen[memberID] = struct{}{}
		ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceGuildLeft{
			CharacterID:    memberID,
			Info:           info,
			RemovedGuildID: removedGuildID,
			RemovedGuild:   removedInfo,
			RemovedMembers: removedMembers,
			Expelled:       expelled,
			Leaving:        leaving,
		})
	}

	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			deliver(m.GetCharacterId(), false)
		}
	}
	for _, m := range removedGuildPb.GetMembers() {
		if m == nil {
			continue
		}
		deliver(m.GetCharacterId(), true)
	}
}

func (ac *AllianceContainer) DisbandAsync(ctx actor.Context, allianceID uint32, guildIDs []uint32) *async.Promise {
	if ac.gs == nil {
		return async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	}
	return ac.gs.guild.RefreshAsync(ctx, guildIDs).Then(func(interface{}) (interface{}, error) {
		ac.Remove(allianceID)
		ac.broadcastDisband(allianceID, guildIDs)
		return nil, nil
	})
}

func (ac *AllianceContainer) broadcastDisband(allianceID uint32, guildIDs []uint32) {
	if ac.gs == nil {
		return
	}
	seen := make(map[uint32]struct{})
	for _, guildID := range guildIDs {
		g := ac.gs.guild.Get(guildID)
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 {
				continue
			}
			if _, dup := seen[memberID]; dup {
				continue
			}
			seen[memberID] = struct{}{}
			ac.deliverToOnlineMember(memberID, &g_actor.DeliverAllianceDisband{
				CharacterID: memberID,
				AllianceID:  allianceID,
			})
		}
	}
}

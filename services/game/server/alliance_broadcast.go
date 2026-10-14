package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	gameconst "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

func (ac *AllianceCache) BroadcastCreate(alliancePb *internal.Alliance) {
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
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceCreate(ch, info, guilds, membershipGuilds)
		})
	}
}

func (ac *AllianceCache) guildInfos(alliancePb *internal.Alliance) []*dto.GuildInfo {
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
		pb := ent.BuildRPC()
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

func (ac *AllianceCache) membershipGuilds(alliancePb *internal.Alliance) []dto.AllianceMembershipChangeGuild {
	if ac.gs == nil || alliancePb == nil {
		return nil
	}
	var out []dto.AllianceMembershipChangeGuild
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := ac.gs.guild.Get(guildID)
		if ent == nil {
			continue
		}
		pb := ent.BuildRPC()
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

func (ac *AllianceCache) recipientIDs(alliancePb *internal.Alliance) []uint32 {
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

func (ac *AllianceCache) forEachOnlineAllianceMember(alliancePb *internal.Alliance, fn func(memberID uint32)) {
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

func (ac *AllianceCache) callOnlineMember(memberID uint32, run func(ch *entity.Character)) {
	if ac.gs == nil || memberID == 0 {
		return
	}
	if ac.gs.characterRuntime == nil || !ac.gs.characterRuntime.Exists(memberID) {
		return
	}
	ac.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
		if ch == nil {
			return
		}
		run(ch)
	})
}

func (ac *AllianceCache) BroadcastNoticeChanged(alliancePb *internal.Alliance, notice string) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceNoticeChanged(ch, info)
		})
	})
	if notice == "" {
		return
	}
	msg := gameconst.AllianceNoticeChangedMessagePrefix + notice
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnMessage(ch, gameconst.MsgPopup, msg)
		})
	})
}

func (ac *AllianceCache) BroadcastLeaderChanged(alliancePb *internal.Alliance, oldLeaderID uint32, newLeaderID uint32) {
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
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceLeaderChanged(ch, allianceID, oldLeaderID, newLeaderID, info, guilds)
		})
		if leaderNotice != "" {
			ac.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
				if ch == nil {
					return
				}
				ch.Listener.OnMessage(ch, gameconst.MsgPinkText, leaderNotice)
			})
		}
	})
}

func (ac *AllianceCache) BroadcastInfoUpdate(alliancePb *internal.Alliance) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceInfoBroadcast(ch, info)
		})
	})
}

func (ac *AllianceCache) broadcastAllianceStateRefresh(alliancePb *internal.Alliance) {
	if ac.gs == nil || alliancePb == nil {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := ac.guildInfos(alliancePb)
	ac.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceMemberRankChanged(ch, info, guilds)
		})
	})
}

func (ac *AllianceCache) allianceCharacterName(alliancePb *internal.Alliance, characterID uint32) string {
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

func (ac *AllianceCache) BroadcastMemberRankChanged(alliancePb *internal.Alliance, characterID uint32, allianceRank uint32) {
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

func (ac *AllianceCache) BroadcastGuildAdded(alliancePb *internal.Alliance, addedGuildPb *internal.Guild) {
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
	var membershipGuild *dto.AllianceMembershipChangeGuild
	if block, ok := entity.AllianceMembershipChangeGuildFromProto(addedGuildPb); ok {
		membershipGuild = &block
	}

	seen := make(map[uint32]struct{})
	send := func(memberID uint32, joining bool) {
		if memberID == 0 {
			return
		}
		if _, dup := seen[memberID]; dup {
			return
		}
		seen[memberID] = struct{}{}
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceGuildAdded(ch, info, guilds, newGuildID, addedInfo, members, joining, membershipGuild)
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
			send(m.GetCharacterId(), false)
		}
	}
	for _, m := range addedGuildPb.GetMembers() {
		if m == nil {
			continue
		}
		send(m.GetCharacterId(), true)
	}
	ac.broadcastAllianceStateRefresh(alliancePb)
}

func (ac *AllianceCache) BroadcastGuildLeft(alliancePb *internal.Alliance, removedGuildPb *internal.Guild, expelled bool) {
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
	send := func(memberID uint32, leaving bool) {
		if memberID == 0 {
			return
		}
		if _, dup := seen[memberID]; dup {
			return
		}
		seen[memberID] = struct{}{}
		ac.callOnlineMember(memberID, func(ch *entity.Character) {
			ch.Listener.OnAllianceGuildLeft(ch, info, removedGuildID, removedInfo, removedMembers, expelled, leaving)
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
			send(m.GetCharacterId(), false)
		}
	}
	for _, m := range removedGuildPb.GetMembers() {
		if m == nil {
			continue
		}
		send(m.GetCharacterId(), true)
	}
}

func (ac *AllianceCache) DisbandAsync(ctx actor.Context, allianceID uint32, guildIDs []uint32) *async.Task {
	if ac.gs == nil {
		return async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	}
	return ac.gs.guild.RefreshAsync(ctx, guildIDs).Do(func() error {
		ac.Remove(allianceID)
		ac.broadcastDisband(allianceID, guildIDs)
		return nil
	})
}

func (ac *AllianceCache) broadcastDisband(allianceID uint32, guildIDs []uint32) {
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
			ac.callOnlineMember(memberID, func(ch *entity.Character) {
				ch.Listener.OnAllianceDisband(ch, allianceID)
			})
		}
	}
}

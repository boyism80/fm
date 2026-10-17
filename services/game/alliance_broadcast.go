package game

import (
	"github.com/asynkron/protoactor-go/actor"
	gameconst "github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

func (s *AllianceSystem) BroadcastCreate(alliancePb *internal.Alliance) {
	if s.gw == nil || alliancePb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := s.guildInfos(alliancePb)
	membershipGuilds := s.membershipGuilds(alliancePb)
	s.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceCreate(ch, info, guilds, membershipGuilds)
		})
	})
}

func (s *AllianceSystem) guildInfos(alliancePb *internal.Alliance) []*dto.GuildInfo {
	var out []*dto.GuildInfo
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
		if ent == nil {
			continue
		}
		if info := GuildInfoFromProto(ent.BuildRPC()); info != nil {
			out = append(out, info)
		}
	}
	return out
}

func (s *AllianceSystem) membershipGuilds(alliancePb *internal.Alliance) []dto.AllianceMembershipChangeGuild {
	if s.gw == nil || alliancePb == nil {
		return nil
	}
	var out []dto.AllianceMembershipChangeGuild
	for _, guildID := range alliancePb.GetGuildIds() {
		ent := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
		if ent == nil {
			continue
		}
		pb := ent.BuildRPC()
		if pb == nil {
			continue
		}
		block, ok := AllianceMembershipChangeGuildFromProto(pb)
		if ok {
			out = append(out, block)
		}
	}
	return out
}

func (s *AllianceSystem) forEachOnlineAllianceMember(alliancePb *internal.Alliance, fn func(memberID uint32)) {
	if s.gw == nil || alliancePb == nil || fn == nil {
		return
	}
	seen := make(map[uint32]struct{})
	for _, guildID := range alliancePb.GetGuildIds() {
		g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
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
			if !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
				continue
			}
			fn(memberID)
		}
	}
}

func (s *AllianceSystem) callOnlineMember(memberID uint32, run func(ch *Character)) {
	if s.gw == nil || memberID == 0 {
		return
	}
	if !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
		return
	}
	s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
		if ch == nil {
			return
		}
		run(ch)
	})
}

func (s *AllianceSystem) BroadcastNoticeChanged(alliancePb *internal.Alliance, notice string) {
	if s.gw == nil || alliancePb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	s.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceNoticeChanged(ch, info)
		})
	})
	if notice == "" {
		return
	}
	msg := gameconst.AllianceNoticeChangedMessagePrefix + notice
	s.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnMessage(ch, gameconst.MsgPopup, msg)
		})
	})
}

func (s *AllianceSystem) BroadcastLeaderChanged(alliancePb *internal.Alliance, oldLeaderID uint32, newLeaderID uint32) {
	if s.gw == nil || alliancePb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	allianceID := alliancePb.GetAllianceId()
	guilds := s.guildInfos(alliancePb)
	leaderNotice := ""
	if name := s.allianceCharacterName(alliancePb, newLeaderID); name != "" {
		leaderNotice = name + gameconst.AllianceLeaderChangedMessageSuffix
	}
	s.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceLeaderChanged(ch, allianceID, oldLeaderID, newLeaderID, info, guilds)
		})
		if leaderNotice != "" {
			s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
				if ch == nil {
					return
				}
				ch.listener.OnMessage(ch, gameconst.MsgPinkText, leaderNotice)
			})
		}
	})
}

func (s *AllianceSystem) BroadcastInfoUpdate(alliancePb *internal.Alliance) {
	if s.gw == nil || alliancePb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	s.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceInfoBroadcast(ch, info)
		})
	})
}

func (s *AllianceSystem) broadcastAllianceStateRefresh(alliancePb *internal.Alliance) {
	if s.gw == nil || alliancePb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	guilds := s.guildInfos(alliancePb)
	s.forEachOnlineAllianceMember(alliancePb, func(memberID uint32) {
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceMemberRankChanged(ch, info, guilds)
		})
	})
}

func (s *AllianceSystem) allianceCharacterName(alliancePb *internal.Alliance, characterID uint32) string {
	if alliancePb == nil {
		return ""
	}
	for _, guildID := range alliancePb.GetGuildIds() {
		g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
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

func (s *AllianceSystem) BroadcastMemberRankChanged(alliancePb *internal.Alliance, characterID uint32, allianceRank uint32) {
	if s.gw == nil || alliancePb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	if info == nil {
		return
	}
	if allianceRank == 0 {
		return
	}
	s.broadcastAllianceStateRefresh(alliancePb)
}

func (s *AllianceSystem) BroadcastGuildAdded(alliancePb *internal.Alliance, addedGuildPb *internal.Guild) {
	if s.gw == nil || alliancePb == nil || addedGuildPb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	addedInfo := GuildInfoFromProto(addedGuildPb)
	if info == nil || addedInfo == nil {
		return
	}
	guilds := s.guildInfos(alliancePb)
	newGuildID := addedGuildPb.GetGuildId()
	members := AllianceGuildMemberRanksFromProto(addedGuildPb)
	var membershipGuild *dto.AllianceMembershipChangeGuild
	if block, ok := AllianceMembershipChangeGuildFromProto(addedGuildPb); ok {
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
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceGuildAdded(ch, info, guilds, newGuildID, addedInfo, members, joining, membershipGuild)
		})
	}

	for _, guildID := range alliancePb.GetGuildIds() {
		g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
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
	s.broadcastAllianceStateRefresh(alliancePb)
}

func (s *AllianceSystem) BroadcastGuildLeft(alliancePb *internal.Alliance, removedGuildPb *internal.Guild, expelled bool) {
	if s.gw == nil || alliancePb == nil || removedGuildPb == nil {
		return
	}
	info := AllianceInfoFromProto(alliancePb)
	removedInfo := GuildInfoFromProto(removedGuildPb)
	if info == nil || removedInfo == nil {
		return
	}
	removedGuildID := removedGuildPb.GetGuildId()
	removedMembers := AllianceGuildMemberRanksFromProto(removedGuildPb)

	seen := make(map[uint32]struct{})
	send := func(memberID uint32, leaving bool) {
		if memberID == 0 {
			return
		}
		if _, dup := seen[memberID]; dup {
			return
		}
		seen[memberID] = struct{}{}
		s.callOnlineMember(memberID, func(ch *Character) {
			ch.listener.OnAllianceGuildLeft(ch, info, removedGuildID, removedInfo, removedMembers, expelled, leaving)
		})
	}

	for _, guildID := range alliancePb.GetGuildIds() {
		g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
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

func (s *AllianceSystem) Disband(allianceID uint32, guildIDs []uint32, guildPbs []*internal.Guild) {
	if s.Get(allianceID) != nil {
		for _, g := range guildPbs {
			s.gw.Systems().Get[*GuildSystem]().store(g)
		}
	}
	s.Remove(allianceID)
	s.broadcastDisband(allianceID, guildIDs)
}

func (s *AllianceSystem) broadcastDisband(allianceID uint32, guildIDs []uint32) {
	if s.gw == nil {
		return
	}
	seen := make(map[uint32]struct{})
	for _, guildID := range guildIDs {
		g := s.gw.Systems().Get[*GuildSystem]().Get(guildID)
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
			s.callOnlineMember(memberID, func(ch *Character) {
				ch.listener.OnAllianceDisband(ch, allianceID)
			})
		}
	}
}

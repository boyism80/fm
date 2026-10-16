package game

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
)

func (s *GuildSystem) SyncGuildEventAsync(ctx actor.Context, evt GuildEventEnvelope, after func(guildID uint32)) *async.Task {
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if s == nil {
		return promise
	}
	return s.UpdateAsync(ctx, evt).Do(func() error {
		if after != nil {
			after(evt.GuildID)
		}
		return nil
	})
}

func (s *GuildSystem) forEachOnlineGuildMember(g *Guild, skipCharacterID uint32, fn func(memberID uint32)) {
	if g == nil || fn == nil {
		return
	}
	for _, m := range g.GetMembers() {
		if m == nil {
			continue
		}
		memberID := m.GetCharacterId()
		if memberID == 0 || memberID == skipCharacterID {
			continue
		}
		if s.gw == nil || !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
			continue
		}
		fn(memberID)
	}
}

func (s *GuildSystem) BroadcastNoticeChanged(guildID uint32) {
	g := s.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	notice := g.notice
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildNoticeChange(ch, gid, notice)
		})
	})
}

func (s *GuildSystem) BroadcastMemberRankChanged(guildID uint32, targetCharacterID uint32) {
	if targetCharacterID == 0 {
		return
	}
	g := s.Get(guildID)
	if g == nil {
		return
	}
	targetRank := uint8(0)
	found := false
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != targetCharacterID {
			continue
		}
		targetRank = m.RankWire()
		found = true
		break
	}
	if !found {
		return
	}
	gid := g.GetGuildId()
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildMemberRankChange(ch, gid, targetCharacterID, targetRank)
		})
	})
}

func (s *GuildSystem) NotifyMemberFieldsChanged(ch *Character) {
	if s == nil || s.gw == nil || ch == nil {
		return
	}
	guildID, inGuild := ch.guild.ID()
	if !inGuild {
		return
	}
	subjectID := ch.GetID()
	level := uint32(ch.GetLevel())
	classID := uint32(ch.class)
	s.mu.Lock()
	g := s.guilds[guildID]
	if g == nil {
		s.mu.Unlock()
		return
	}
	if m := g.FindMember(subjectID); m != nil {
		m.level = level
		m.classID = classID
	}
	s.mu.Unlock()

	s.BroadcastMemberFieldsChanged(guildID, subjectID, level, classID)
}

func (s *GuildSystem) BroadcastMemberFieldsChanged(guildID uint32, subjectID uint32, level uint32, classID uint32) {
	if subjectID == 0 {
		return
	}
	g := s.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildMemberFieldsChange(ch, gid, subjectID, level, classID)
		})
	})
}

func (s *GuildSystem) BroadcastMemberOnlineChanged(guildID uint32, subjectID uint32, online bool) {
	if subjectID == 0 {
		return
	}
	g := s.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	s.forEachOnlineGuildMember(g, subjectID, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildMemberOnlineChange(ch, gid, subjectID, online)
		})
	})
	allianceID, inAlliance := g.GetAllianceID()
	if inAlliance && allianceID > 0 {
		s.gw.Systems().Get[*AllianceSystem]().BroadcastMemberOnlineChanged(allianceID, gid, subjectID, online)
	}
}

func (s *GuildSystem) BroadcastEmblemChanged(guildID uint32) {
	g := s.Get(guildID)
	if g == nil || g.logo == nil {
		return
	}
	gid := g.GetGuildId()
	logoBG := uint16(g.logo.logoBG)
	logoBGColor := uint8(g.logo.logoBGColor)
	logo := uint16(g.logo.logo)
	logoColor := uint8(g.logo.logoColor)
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildEmblemChange(ch, gid, logoBG, logoBGColor, logo, logoColor)
		})
	})
}

func (s *GuildSystem) BroadcastCapacityChanged(guildID uint32) {
	g := s.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	capacity := uint8(g.capacity)
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildCapacityChange(ch, gid, capacity)
		})
	})
}

func (s *GuildSystem) BroadcastGPChanged(guildID uint32, amount int32) {
	g := s.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gp := g.gp
	level := g.Level()
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildGPChange(ch, gid, gp, level, amount)
		})
	})
}

func (s *GuildSystem) BroadcastRankTitlesChanged(guildID uint32) {
	g := s.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	rankTitles := g.rankTitles
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildRankTitleChange(ch, gid, rankTitles)
		})
	})
}

func (s *GuildSystem) BroadcastMemberJoined(guildID uint32, joinerCharacterID uint32) {
	if joinerCharacterID == 0 {
		return
	}
	g := s.Get(guildID)
	if g == nil {
		return
	}
	var joinerStatus dto.GuildMemberStatus
	found := false
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != joinerCharacterID {
			continue
		}
		joinerStatus = m.BuildS2C()
		found = true
		break
	}
	if !found {
		return
	}
	gid := g.GetGuildId()
	s.forEachOnlineGuildMember(g, joinerCharacterID, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildNewMember(ch, gid, joinerStatus)
		})
	})
}

func (s *GuildSystem) BroadcastMemberLeft(prevGuild *Guild, leftCharacterID uint32, expelled bool) {
	if prevGuild == nil || leftCharacterID == 0 {
		return
	}
	leftName := ""
	for _, m := range prevGuild.GetMembers() {
		if m == nil || m.GetCharacterId() != leftCharacterID {
			continue
		}
		leftName = m.GetCharacterName()
		break
	}
	if leftName == "" {
		return
	}
	guildID := prevGuild.GetGuildId()
	for _, m := range prevGuild.GetMembers() {
		if m == nil {
			continue
		}
		memberID := m.GetCharacterId()
		if memberID == 0 {
			continue
		}
		if memberID == leftCharacterID {
			if expelled {
				s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
					if ch == nil {
						return
					}
					ch.listener.OnGuildExpelledSelf(ch, guildID)
				})
			} else {
				s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
					if ch == nil {
						return
					}
					ch.listener.OnGuildLeaveSelf(ch)
				})
			}
			continue
		}
		if s.gw == nil || !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildMemberLeft(ch, guildID, leftCharacterID, leftName, expelled)
		})
	}
}

func (s *GuildSystem) BroadcastDisbanded(prevGuild *Guild, memberCharacterIDs []uint32) {
	if prevGuild == nil {
		return
	}
	guildID := prevGuild.GetGuildId()
	memberIDs := memberCharacterIDs
	if len(memberIDs) == 0 {
		memberIDs = prevGuild.MemberCharacterIDs()
	}
	for _, memberID := range memberIDs {
		if memberID == 0 {
			continue
		}
		if s.gw == nil || !s.gw.Systems().Get[*DispatchSystem]().Exists(memberID) {
			continue
		}
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnGuildDisbandSelf(ch, guildID)
		})
	}
}

func (s *GuildSystem) BroadcastMessage(guildID uint32, messageType constant.ServerMessageType, message string) {
	g := s.Get(guildID)
	if g == nil {
		return
	}
	s.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnMessage(ch, messageType, message)
		})
	})
}

func (s *GuildSystem) BroadcastMultiChat(guildID uint32, senderCharacterID uint32, senderName string, message string) {
	if senderCharacterID == 0 {
		return
	}
	g := s.Get(guildID)
	if g == nil {
		return
	}
	mode := pconst.MultiChatModeGuild
	s.forEachOnlineGuildMember(g, senderCharacterID, func(memberID uint32) {
		s.gw.Systems().Get[*DispatchSystem]().Visit(memberID, func(_ actor.Context, ch *Character) {
			if ch == nil {
				return
			}
			ch.listener.OnMultiChat(ch, mode, senderName, message)
		})
	})
}

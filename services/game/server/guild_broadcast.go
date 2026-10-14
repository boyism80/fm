package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

func (gc *GuildCache) SyncGuildEventAsync(ctx actor.Context, evt GuildEventEnvelope, after func(guildID uint32)) *async.Task {
	promise := async.NewTask(ctx, core.InternalRPCPerStepTimeout)
	if gc == nil {
		return promise
	}
	return gc.UpdateAsync(ctx, evt).Do(func() error {
		if after != nil {
			after(evt.GuildID)
		}
		return nil
	})
}

func (gc *GuildCache) forEachOnlineGuildMember(g *entity.Guild, skipCharacterID uint32, fn func(memberID uint32)) {
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
		if gc.gs == nil || gc.gs.characterRuntime == nil || !gc.gs.characterRuntime.Exists(memberID) {
			continue
		}
		fn(memberID)
	}
}

func (gc *GuildCache) BroadcastNoticeChanged(guildID uint32) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	notice := g.Notice
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildNoticeChange(ch, gid, notice)
		})
	})
}

func (gc *GuildCache) BroadcastMemberRankChanged(guildID uint32, targetCharacterID uint32) {
	if targetCharacterID == 0 {
		return
	}
	g := gc.Get(guildID)
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
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildMemberRankChange(ch, gid, targetCharacterID, targetRank)
		})
	})
}

func (gc *GuildCache) NotifyMemberFieldsChanged(ch *entity.Character) {
	if gc == nil || gc.gs == nil || ch == nil {
		return
	}
	guildID, inGuild := ch.Guild.ID()
	if !inGuild {
		return
	}
	subjectID := ch.GetID()
	level := uint32(ch.GetLevel())
	classID := uint32(ch.Class)
	gc.mu.Lock()
	g := gc.guilds[guildID]
	if g == nil {
		gc.mu.Unlock()
		return
	}
	if m := g.FindMember(subjectID); m != nil {
		m.Level = level
		m.ClassID = classID
	}
	gc.mu.Unlock()

	gc.BroadcastMemberFieldsChanged(guildID, subjectID, level, classID)
}

func (gc *GuildCache) BroadcastMemberFieldsChanged(guildID uint32, subjectID uint32, level uint32, classID uint32) {
	if subjectID == 0 {
		return
	}
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildMemberFieldsChange(ch, gid, subjectID, level, classID)
		})
	})
}

func (gc *GuildCache) BroadcastMemberOnlineChanged(guildID uint32, subjectID uint32, online bool) {
	if subjectID == 0 {
		return
	}
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gc.forEachOnlineGuildMember(g, subjectID, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildMemberOnlineChange(ch, gid, subjectID, online)
		})
	})
	allianceID, inAlliance := g.GetAllianceID()
	if inAlliance && allianceID > 0 {
		gc.gs.alliance.BroadcastMemberOnlineChanged(allianceID, gid, subjectID, online)
	}
}

func (gc *GuildCache) BroadcastEmblemChanged(guildID uint32) {
	g := gc.Get(guildID)
	if g == nil || g.Logo == nil {
		return
	}
	gid := g.GetGuildId()
	logoBG := uint16(g.Logo.LogoBG)
	logoBGColor := uint8(g.Logo.LogoBGColor)
	logo := uint16(g.Logo.Logo)
	logoColor := uint8(g.Logo.LogoColor)
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildEmblemChange(ch, gid, logoBG, logoBGColor, logo, logoColor)
		})
	})
}

func (gc *GuildCache) BroadcastCapacityChanged(guildID uint32) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	capacity := uint8(g.Capacity)
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildCapacityChange(ch, gid, capacity)
		})
	})
}

func (gc *GuildCache) BroadcastGPChanged(guildID uint32, amount int32) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gp := g.GP
	level := g.Level()
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildGPChange(ch, gid, gp, level, amount)
		})
	})
}

func (gc *GuildCache) BroadcastRankTitlesChanged(guildID uint32) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	rankTitles := g.RankTitles
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildRankTitleChange(ch, gid, rankTitles)
		})
	})
}

func (gc *GuildCache) BroadcastMemberJoined(guildID uint32, joinerCharacterID uint32) {
	if joinerCharacterID == 0 {
		return
	}
	g := gc.Get(guildID)
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
	gc.forEachOnlineGuildMember(g, joinerCharacterID, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildNewMember(ch, gid, joinerStatus)
		})
	})
}

func (gc *GuildCache) BroadcastMemberLeft(prevGuild *entity.Guild, leftCharacterID uint32, expelled bool) {
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
				gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
					if ch == nil {
						return
					}
					ch.Listener.OnGuildExpelledSelf(ch, guildID)
				})
			} else {
				gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
					if ch == nil {
						return
					}
					ch.Listener.OnGuildLeaveSelf(ch)
				})
			}
			continue
		}
		if gc.gs == nil || gc.gs.characterRuntime == nil || !gc.gs.characterRuntime.Exists(memberID) {
			continue
		}
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildMemberLeft(ch, guildID, leftCharacterID, leftName, expelled)
		})
	}
}

func (gc *GuildCache) BroadcastDisbanded(prevGuild *entity.Guild, memberCharacterIDs []uint32) {
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
		if gc.gs == nil || gc.gs.characterRuntime == nil || !gc.gs.characterRuntime.Exists(memberID) {
			continue
		}
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnGuildDisbandSelf(ch, guildID)
		})
	}
}

func (gc *GuildCache) BroadcastMessage(guildID uint32, messageType constant.ServerMessageType, message string) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnMessage(ch, messageType, message)
		})
	})
}

func (gc *GuildCache) BroadcastMultiChat(guildID uint32, senderCharacterID uint32, senderName string, message string) {
	if senderCharacterID == 0 {
		return
	}
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	mode := pconst.MultiChatModeGuild
	gc.forEachOnlineGuildMember(g, senderCharacterID, func(memberID uint32) {
		gc.gs.GetDispatchSystem().CallCharacter(memberID, func(_ actor.Context, ch *entity.Character) {
			if ch == nil {
				return
			}
			ch.Listener.OnMultiChat(ch, mode, senderName, message)
		})
	})
}

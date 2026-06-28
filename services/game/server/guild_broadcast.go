package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
)

func (gc *GuildContainer) ApplyEventAsync(ctx actor.Context, evt GuildEventEnvelope, after func(guildID uint32)) *async.Promise {
	promise := async.NewPromise(ctx, core.InternalRPCPerStepTimeout)
	if gc == nil {
		return promise
	}
	return gc.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		if after != nil {
			after(evt.GuildID)
		}
		return nil, nil
	})
}

func (gc *GuildContainer) guildGet(guildID uint32) *entity.Guild {
	if gc == nil || guildID == 0 {
		return nil
	}
	return gc.Get(guildID)
}

func (gc *GuildContainer) isCharacterOnline(characterID uint32) bool {
	if gc == nil || gc.gs == nil || characterID == 0 {
		return false
	}
	if gc.gs.characterRuntime == nil {
		return false
	}
	return gc.gs.characterRuntime.Exists(characterID)
}

func (gc *GuildContainer) forEachOnlineGuildMember(g *entity.Guild, skipCharacterID uint32, fn func(memberID uint32)) {
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
		if !gc.isCharacterOnline(memberID) {
			continue
		}
		fn(memberID)
	}
}

func (gc *GuildContainer) BroadcastNoticeChanged(guildID uint32) {
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	notice := g.Notice
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildNoticeChange{
			CharacterID: memberID,
			GuildID:     gid,
			Notice:      notice,
		})
	})
}

func (gc *GuildContainer) BroadcastMemberRankChanged(guildID uint32, targetCharacterID uint32) {
	if targetCharacterID == 0 {
		return
	}
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	targetRank := uint8(0)
	found := false
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != targetCharacterID {
			continue
		}
		targetRank = entity.GuildMemberRankWire(m.GetRank())
		found = true
		break
	}
	if !found {
		return
	}
	gid := g.GetGuildId()
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberRankChange{
			CharacterID: memberID,
			GuildID:     gid,
			TargetID:    targetCharacterID,
			GuildRank:   targetRank,
		})
	})
}

func (gc *GuildContainer) NotifyMemberFieldsChanged(ch *entity.Character) {
	if gc == nil || gc.gs == nil || ch == nil {
		return
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild {
		return
	}
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	subjectID := ch.GetID()
	level := uint32(ch.GetLevel())
	classID := uint32(ch.Class)
	if m := g.FindMember(subjectID); m != nil {
		m.Level = level
		m.ClassID = classID
	}
	gc.BroadcastMemberFieldsChanged(guildID, subjectID, level, classID)
}

func (gc *GuildContainer) BroadcastMemberFieldsChanged(guildID uint32, subjectID uint32, level uint32, classID uint32) {
	if subjectID == 0 {
		return
	}
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberFieldsChange{
			CharacterID: memberID,
			GuildID:     gid,
			SubjectID:   subjectID,
			Level:       level,
			ClassID:     classID,
		})
	})
}

func (gc *GuildContainer) BroadcastMemberOnlineChanged(guildID uint32, subjectID uint32, online bool) {
	if subjectID == 0 {
		return
	}
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gc.forEachOnlineGuildMember(g, subjectID, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberOnlineChange{
			CharacterID: memberID,
			GuildID:     gid,
			SubjectID:   subjectID,
			Online:      online,
		})
	})
	allianceID, inAlliance := g.GetAllianceID()
	if inAlliance && allianceID > 0 {
		gc.gs.alliance.DeliverMemberOnlineChange(allianceID, gid, subjectID, online)
	}
}

func (gc *GuildContainer) BroadcastEmblemChanged(guildID uint32) {
	g := gc.guildGet(guildID)
	if g == nil || g.Logo == nil {
		return
	}
	gid := g.GetGuildId()
	logo := g.Logo
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildEmblemChange{
			CharacterID: memberID,
			GuildID:     gid,
			LogoBG:      uint16(logo.LogoBG),
			LogoBGColor: uint8(logo.LogoBGColor),
			Logo:        uint16(logo.Logo),
			LogoColor:   uint8(logo.LogoColor),
		})
	})
}

func (gc *GuildContainer) BroadcastCapacityChanged(guildID uint32) {
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	capacity := uint8(g.Capacity)
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildCapacityChange{
			CharacterID: memberID,
			GuildID:     gid,
			Capacity:    capacity,
		})
	})
}

func (gc *GuildContainer) BroadcastRankTitlesChanged(guildID uint32) {
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	rankTitles := g.RankTitles
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildRankTitleChange{
			CharacterID: memberID,
			GuildID:     gid,
			RankTitles:  rankTitles,
		})
	})
}

func (gc *GuildContainer) BroadcastMemberJoined(guildID uint32, joinerCharacterID uint32) {
	if joinerCharacterID == 0 {
		return
	}
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	var joinerStatus dto.GuildMemberStatus
	found := false
	for _, m := range g.GetMembers() {
		if m == nil || m.GetCharacterId() != joinerCharacterID {
			continue
		}
		joinerStatus = entity.GuildMemberToDTO(m)
		found = true
		break
	}
	if !found {
		return
	}
	gid := g.GetGuildId()
	gc.forEachOnlineGuildMember(g, joinerCharacterID, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildNewMember{
			CharacterID: memberID,
			GuildID:     gid,
			Member:      joinerStatus,
		})
	})
}

func (gc *GuildContainer) BroadcastMemberLeft(prevGuild *entity.Guild, leftCharacterID uint32, expelled bool) {
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
				gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildExpelSelf{
					CharacterID: memberID,
					GuildID:     guildID,
				})
			} else {
				gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildLeaveSelf{
					CharacterID: memberID,
				})
			}
			continue
		}
		if !gc.isCharacterOnline(memberID) {
			continue
		}
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberLeft{
			CharacterID: memberID,
			GuildID:     guildID,
			TargetID:    leftCharacterID,
			TargetName:  leftName,
			WasExpelled: expelled,
		})
	}
}

func (gc *GuildContainer) BroadcastDisbanded(prevGuild *entity.Guild, memberCharacterIDs []uint32) {
	if prevGuild == nil {
		return
	}
	guildID := prevGuild.GetGuildId()
	memberIDs := memberCharacterIDs
	if len(memberIDs) == 0 {
		memberIDs = entity.GuildMemberCharacterIDs(prevGuild.GetMembers())
	}
	for _, memberID := range memberIDs {
		if memberID == 0 {
			continue
		}
		if !gc.isCharacterOnline(memberID) {
			continue
		}
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildDisbandSelf{
			CharacterID: memberID,
			GuildID:     guildID,
		})
	}
}

func (gc *GuildContainer) BroadcastMultiChat(guildID uint32, senderCharacterID uint32, senderName string, message string) {
	if senderCharacterID == 0 {
		return
	}
	g := gc.guildGet(guildID)
	if g == nil {
		return
	}
	mode := pconst.MultiChatModeGuild
	gc.forEachOnlineGuildMember(g, senderCharacterID, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverMultiChat{
			CharacterID: memberID,
			Mode:        mode,
			SenderName:  senderName,
			Message:     message,
		})
	})
}

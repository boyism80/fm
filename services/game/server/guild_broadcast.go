package server

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	g_actor "github.com/boyism80/fm/services/game/actor"
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
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildNoticeChange{
			CharacterID: memberID,
			GuildID:     gid,
			Notice:      notice,
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

func (gc *GuildCache) NotifyMemberFieldsChanged(ch *entity.Character) {
	if gc == nil || gc.gs == nil || ch == nil {
		return
	}
	guildID, inGuild := ch.Guild.ID()
	if !inGuild {
		return
	}
	g := gc.Get(guildID)
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
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberFieldsChange{
			CharacterID: memberID,
			GuildID:     gid,
			SubjectID:   subjectID,
			Level:       level,
			ClassID:     classID,
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

func (gc *GuildCache) BroadcastEmblemChanged(guildID uint32) {
	g := gc.Get(guildID)
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

func (gc *GuildCache) BroadcastCapacityChanged(guildID uint32) {
	g := gc.Get(guildID)
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

func (gc *GuildCache) BroadcastGPChanged(guildID uint32, amount int32) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gid := g.GetGuildId()
	gp := g.GP
	level := g.Level()
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildGPChange{
			CharacterID: memberID,
			GuildID:     gid,
			GP:          gp,
			Level:       level,
			Amount:      amount,
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
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildRankTitleChange{
			CharacterID: memberID,
			GuildID:     gid,
			RankTitles:  rankTitles,
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
		if gc.gs == nil || gc.gs.characterRuntime == nil || !gc.gs.characterRuntime.Exists(memberID) {
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

func (gc *GuildCache) BroadcastDisbanded(prevGuild *entity.Guild, memberCharacterIDs []uint32) {
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
		if gc.gs == nil || gc.gs.characterRuntime == nil || !gc.gs.characterRuntime.Exists(memberID) {
			continue
		}
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildDisbandSelf{
			CharacterID: memberID,
			GuildID:     guildID,
		})
	}
}

func (gc *GuildCache) BroadcastMessage(guildID uint32, messageType constant.ServerMessageType, message string) {
	g := gc.Get(guildID)
	if g == nil {
		return
	}
	gc.forEachOnlineGuildMember(g, 0, func(memberID uint32) {
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverMessage{
			CharacterID: memberID,
			MessageType: messageType,
			Message:     message,
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
		gc.gs.EnsureSend(nil, memberID, &g_actor.DeliverMultiChat{
			CharacterID: memberID,
			Mode:        mode,
			SenderName:  senderName,
			Message:     message,
		})
	})
}

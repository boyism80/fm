package server

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type CharacterListenerImpl struct {
	gs *GameServer
}

func (l *CharacterListenerImpl) OnDialog(ch *game.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.Dialog{
		NPC:  npc,
		Type: constant.DialogTypeDefault,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogYesNo(ch *game.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogAccept(ch *game.Character, npc uint32, message string, enableEscape bool) {
	dialogPacket := &response.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogList(ch *game.Character, npc uint32, message string, selections []string) {
	dialogPacket := &response.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogStyle(ch *game.Character, npc uint32, message string, styles []uint32) {
	ch.Send(&response.DialogStyle{
		NPC:    npc,
		Text:   message,
		Styles: styles,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyCapacity(ch *game.Character, capacity uint8) {
	ch.Send(&response.BuddyCapacityUpdate{Capacity: capacity}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogInput(ch *game.Character, npc uint32, message string) {
	dialogPacket := &response.DialogInput{
		NPC:  npc,
		Text: message,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogNumber(ch *game.Character, npc uint32, message string, def, min, max int32) {
	ch.Send(&response.DialogNumber{
		NPC:     npc,
		Text:    message,
		Default: def,
		Min:     min,
		Max:     max,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnChat(ch *game.Character, message string, highlight bool, dontRecordHistory bool) {
	chatPacket := &response.NormalChat{
		CharacterId:       ch.GetID(),
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(chatPacket, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnMesoChanged(ch *game.Character, meso int32) {
	ch.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnKeyMap(ch *game.Character) {
	ch.Send(&response.KeyMap{Slots: ch.KeyLayout().Bindings()}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSkillMacros(ch *game.Character) {
	macros := ch.Skills().Macros().Entries()
	entries := make([]response.SkillMacroEntry, 0, len(macros))
	for _, macro := range macros {
		entries = append(entries, response.SkillMacroEntry{Name: macro.Name, Shout: macro.Shout, Skills: macro.Skills})
	}
	ch.Send(&response.SkillMacro{Entries: entries}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnChalkboard(ch *game.Character) {
	ch.Broadcast(&response.Chalkboard{CharacterID: ch.GetID(), Text: ch.Chalkboard()}, nil)
}

func (l *CharacterListenerImpl) OnFaceExpression(ch *game.Character, expression uint32) {
	ch.Broadcast(&response.FaceExpression{CharacterID: ch.GetID(), Expression: expression}, nil)
}

func (l *CharacterListenerImpl) OnItemEffect(ch *game.Character) {
	ch.Broadcast(&response.ItemEffect{CharacterID: ch.GetID(), ItemID: ch.ItemEffectID()}, nil)
}

func (l *CharacterListenerImpl) OnChargeSkillEffect(ch *game.Character, skillID uint32, level, flags, speed uint8) {
	ch.Broadcast(&response.ChargeSkillEffect{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		Level:       level,
		Flags:       flags,
		Speed:       speed,
	}, nil)
}

func (l *CharacterListenerImpl) OnChargeSkillCancelled(ch *game.Character, skillID uint32) {
	ch.Broadcast(&response.CancelChargeSkillEffect{CharacterID: ch.GetID(), SkillID: skillID}, nil)
}

func (l *CharacterListenerImpl) OnMessage(ch *game.Character, messageType constant.ServerMessageType, message string) {
	l.OnNotice(ch, messageType, message, 0, false)
}

func (l *CharacterListenerImpl) OnNotice(ch *game.Character, messageType constant.ServerMessageType, message string, channel int, ear bool) {
	ch.Send(&response.Notice{
		Message: message,
		Type:    messageType,
		Channel: channel,
		MegaEar: ear,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMegaphone(ch *game.Character, megaphone *response.Megaphone) {
	ch.Send(megaphone, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessageBoxFailed(ch *game.Character) {
	ch.Send(&response.MessageBoxFailed{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) Disconnect(ch *game.Character) {
	c, ok := ch.Client().(*client.GameClient)
	if !ok {
		return
	}
	_ = c.GetConnection().Close()
}

func (l *CharacterListenerImpl) OnClock(ch *game.Character, seconds int32) {
	ch.Send(&response.Clock{
		Seconds: seconds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStopClock(ch *game.Character) {
	ch.Send(&response.StopClock{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyCreated(ch *game.Character, partyID uint32) {
	ch.Send(&response.PartyCreated{
		PartyID: partyID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowGuildInfo(ch *game.Character) {
	if l.gs == nil {
		return
	}
	guildID, ok := ch.Guild().ID()
	if !ok {
		return
	}
	ent := l.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if ent == nil {
		return
	}
	info := ent.BuildS2C()
	if info == nil {
		return
	}
	ch.Send(&response.GuildShowInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowAllianceInfo(ch *game.Character) {
	if l.gs == nil {
		return
	}
	allianceID := uint32(0)
	if guildID, ok := ch.Guild().ID(); ok {
		if g := l.gs.systems.Get[*game.GuildSystem]().Get(guildID); g != nil {
			if id, inAlliance := g.GetAllianceID(); inAlliance {
				allianceID = id
			}
		}
	}
	info, guilds := l.gs.systems.Get[*game.AllianceSystem]().ShowData(allianceID)
	_ = ch.Send(&response.AllianceShowInfo{Info: info}, types.SEND_POLICY_ENCRYPT)
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceShowGuilds{Guilds: guilds}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceUpdateInfo{Info: info}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceUpdateInfo(ch *game.Character) {
	if l.gs == nil {
		return
	}
	allianceID := uint32(0)
	if guildID, ok := ch.Guild().ID(); ok {
		if g := l.gs.systems.Get[*game.GuildSystem]().Get(guildID); g != nil {
			if id, inAlliance := g.GetAllianceID(); inAlliance {
				allianceID = id
			}
		}
	}
	if allianceID == 0 {
		return
	}
	info, _ := l.gs.systems.Get[*game.AllianceSystem]().ShowData(allianceID)
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateInfo{Info: info}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBroadcastGuildAppearance(ch *game.Character) {
	guildID, ok := ch.Guild().ID()
	if !ok || l.gs == nil {
		return
	}
	ent := l.gs.systems.Get[*game.GuildSystem]().Get(guildID)
	if ent == nil {
		return
	}
	characterID := ch.GetID()
	broadcastOpt := &game.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   ent.Name(),
	}, broadcastOpt)
	iconPkt := &response.LoadGuildIcon{
		CharacterID: characterID,
	}
	if ent.Logo() != nil {
		iconPkt.HasGuildIcon = true
		iconPkt.LogoBG = uint16(ent.Logo().LogoBG())
		iconPkt.LogoBGColor = uint8(ent.Logo().LogoBGColor())
		iconPkt.Logo = uint16(ent.Logo().Logo())
		iconPkt.LogoColor = uint8(ent.Logo().LogoColor())
	}
	ch.Broadcast(iconPkt, broadcastOpt)
}

func (l *CharacterListenerImpl) OnPartyInvite(ch *game.Character, partyID uint32, inviterName string, partySearch bool) {
	ch.Send(&response.PartyInvite{
		PartyID:     partyID,
		InviterName: inviterName,
		PartySearch: partySearch,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildInvite(ch *game.Character, guildID uint32, inviterName string) {
	_ = ch.Send(&response.GuildInvite{
		GuildID:     guildID,
		InviterName: inviterName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceInvite(ch *game.Character, inviterGuildID uint32, inviterName string, allianceName string) {
	_ = ch.Send(&response.AllianceInvite{
		InviterGuildID: inviterGuildID,
		InviterName:    inviterName,
		AllianceName:   allianceName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceCreate(ch *game.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, membershipGuilds []dto.AllianceMembershipChangeGuild) {
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceCreate{
		Info:   info,
		Guilds: guilds,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceShowInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceShowGuilds{
		Guilds: guilds,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceChangeMembership{
		InAlliance: true,
		AllianceID: info.AllianceID,
		Guilds:     membershipGuilds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceDisband(ch *game.Character, allianceID uint32) {
	_ = ch.Send(&response.AllianceDisband{
		AllianceID: allianceID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceInfoBroadcast(ch *game.Character, info *dto.AllianceInfo) {
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceNoticeChanged(ch *game.Character, info *dto.AllianceInfo) {
	l.OnAllianceInfoBroadcast(ch, info)
}

func (l *CharacterListenerImpl) OnAllianceLeaderChanged(ch *game.Character, allianceID uint32, oldLeaderID uint32, newLeaderID uint32, info *dto.AllianceInfo, guilds []*dto.GuildInfo) {
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceChangeLeader{
		AllianceID:  allianceID,
		OldLeaderID: oldLeaderID,
		NewLeaderID: newLeaderID,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceUpdateLeader{
		AllianceID:  allianceID,
		OldLeaderID: oldLeaderID,
		NewLeaderID: newLeaderID,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceUpdateInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceShowGuilds{
		Guilds: guilds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberRankChanged(ch *game.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo) {
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceShowGuilds{
		Guilds: guilds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceGuildAdded(ch *game.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, newGuildID uint32, addedGuild *dto.GuildInfo, members []dto.AllianceGuildMemberRank, joining bool, membershipGuild *dto.AllianceMembershipChangeGuild) {
	if info == nil || addedGuild == nil {
		return
	}
	if joining {
		_ = ch.Send(&response.AllianceShowInfo{
			Info: info,
		}, types.SEND_POLICY_ENCRYPT)
		_ = ch.Send(&response.AllianceShowGuilds{
			Guilds: guilds,
		}, types.SEND_POLICY_ENCRYPT)
		if membershipGuild != nil {
			_ = ch.Send(&response.AllianceChangeMembership{
				InAlliance: true,
				AllianceID: info.AllianceID,
				Guilds:     []dto.AllianceMembershipChangeGuild{*membershipGuild},
			}, types.SEND_POLICY_ENCRYPT)
		}
		return
	}
	_ = ch.Send(&response.AllianceAddGuild{
		Info:       info,
		NewGuildID: newGuildID,
		Guild:      addedGuild,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceChangeGuildMembers{
		Added:      true,
		AllianceID: info.AllianceID,
		GuildID:    newGuildID,
		Members:    members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceGuildLeft(ch *game.Character, info *dto.AllianceInfo, removedGuildID uint32, removedGuild *dto.GuildInfo, removedMembers []dto.AllianceGuildMemberRank, expelled bool, leaving bool) {
	if info == nil || removedGuild == nil {
		return
	}
	_ = ch.Send(&response.AllianceRemoveGuild{
		Info:           info,
		RemovedGuildID: removedGuildID,
		Guild:          removedGuild,
		Expelled:       expelled,
	}, types.SEND_POLICY_ENCRYPT)
	if leaving {
		return
	}
	_ = ch.Send(&response.AllianceChangeGuildMembers{
		Added:      false,
		AllianceID: 0,
		GuildID:    removedGuildID,
		Members:    removedMembers,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildNewMember(ch *game.Character, guildID uint32, member dto.GuildMemberStatus) {
	_ = ch.Send(&response.GuildNewMember{
		GuildID: guildID,
		Member:  member,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildLeaveSelf(ch *game.Character) {
	ch.Guild().SetID(nil)
	_ = ch.Send(&response.GuildShowInfo{
		Info: nil,
	}, types.SEND_POLICY_ENCRYPT)
	characterID := ch.GetID()
	broadcastOpt := &game.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   "",
	}, broadcastOpt)
	ch.Broadcast(&response.LoadGuildIcon{
		CharacterID: characterID,
	}, broadcastOpt)
}

func (l *CharacterListenerImpl) OnGuildExpelledSelf(ch *game.Character, guildID uint32) {
	ch.Guild().SetID(nil)
	_ = ch.Send(&response.GuildMemberLeft{
		GuildID:     guildID,
		CharacterID: ch.GetID(),
		Name:        ch.GetName(),
		WasExpelled: true,
	}, types.SEND_POLICY_ENCRYPT)
	characterID := ch.GetID()
	broadcastOpt := &game.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   "",
	}, broadcastOpt)
	ch.Broadcast(&response.LoadGuildIcon{
		CharacterID: characterID,
	}, broadcastOpt)
}

func (l *CharacterListenerImpl) OnGuildMemberLeft(ch *game.Character, guildID uint32, targetID uint32, targetName string, wasExpelled bool) {
	_ = ch.Send(&response.GuildMemberLeft{
		GuildID:     guildID,
		CharacterID: targetID,
		Name:        targetName,
		WasExpelled: wasExpelled,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildRankTitleChange(ch *game.Character, guildID uint32, rankTitles [5]string) {
	_ = ch.Send(&response.GuildRankTitleChange{
		GuildID:    guildID,
		RankTitles: rankTitles,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberRankChange(ch *game.Character, guildID uint32, targetID uint32, guildRank uint8) {
	_ = ch.Send(&response.GuildChangeRank{
		GuildID:     guildID,
		CharacterID: targetID,
		GuildRank:   guildRank,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildEmblemChange(ch *game.Character, guildID uint32, logoBG uint16, logoBGColor uint8, logo uint16, logoColor uint8) {
	_ = ch.Send(&response.GuildEmblemChange{
		GuildID:     guildID,
		LogoBG:      logoBG,
		LogoBGColor: logoBGColor,
		Logo:        logo,
		LogoColor:   logoColor,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnBroadcastGuildAppearance(ch)
}

func (l *CharacterListenerImpl) OnGuildNoticeChange(ch *game.Character, guildID uint32, notice string) {
	_ = ch.Send(&response.GuildNotice{
		GuildID: guildID,
		Notice:  notice,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildCapacityChange(ch *game.Character, guildID uint32, capacity uint8) {
	_ = ch.Send(&response.GuildCapacityChange{
		GuildID:  guildID,
		Capacity: capacity,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildGPChange(ch *game.Character, guildID uint32, gp uint32, level uint32, amount int32) {
	_ = ch.Send(&response.GuildUpdateGP{
		GuildID:    guildID,
		GP:         gp,
		GuildLevel: level,
	}, types.SEND_POLICY_ENCRYPT)
	if amount != 0 {
		_ = ch.Send(&response.ShowGPGain{Amount: amount}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnGuildRanking(ch *game.Character, npcID uint32, entries []dto.GuildRankingEntry) {
	_ = ch.Send(&response.GuildShowRanks{
		NPCID:   npcID,
		Entries: entries,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberOnlineChange(ch *game.Character, guildID uint32, subjectCharacterID uint32, online bool) {
	_ = ch.Send(&response.GuildMemberOnline{
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Online:      online,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberFieldsChange(ch *game.Character, guildID uint32, subjectCharacterID uint32, level uint32, classID uint32) {
	_ = ch.Send(&response.GuildMemberLevelClassUpdate{
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Level:       level,
		ClassID:     classID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberOnlineChange(ch *game.Character, allianceID uint32, guildID uint32, subjectCharacterID uint32, online bool) {
	_ = ch.Send(&response.AllianceMemberOnline{
		AllianceID:  allianceID,
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Online:      online,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberFieldsChange(ch *game.Character, allianceID uint32, guildID uint32, subjectCharacterID uint32, level uint32, classID uint32) {
	_ = ch.Send(&response.AllianceUpdateMember{
		AllianceID:  allianceID,
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Level:       level,
		ClassID:     classID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildDisbandSelf(ch *game.Character, guildID uint32) {
	if gid, ok := ch.Guild().ID(); ok && gid == guildID {
		_ = ch.Send(&response.GuildDisband{
			GuildID: guildID,
		}, types.SEND_POLICY_ENCRYPT)
	}
	ch.Guild().SetID(nil)
	_ = ch.Send(&response.GuildShowInfo{
		Info: nil,
	}, types.SEND_POLICY_ENCRYPT)
	characterID := ch.GetID()
	broadcastOpt := &game.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   "",
	}, broadcastOpt)
	ch.Broadcast(&response.LoadGuildIcon{
		CharacterID: characterID,
	}, broadcastOpt)
}

func (l *CharacterListenerImpl) OnMultiChat(ch *game.Character, mode pconst.MultiChatMode, senderName string, message string) {
	_ = ch.Send(&response.MultiChat{
		Mode:    mode,
		Name:    senderName,
		Message: message,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyStatusMessage(ch *game.Character, code pconst.PartyStatusCode, name string) {
	ch.Send(&response.PartyStatusMessage{
		Code: code,
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateJoin(ch *game.Character, forChannel int32, partyID uint32, joinName string, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateJoin{
		ForChannel:           forChannel,
		PartyID:              partyID,
		JoiningCharacterName: joinName,
		LeaderCharacterID:    leaderID,
		Members:              members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateExpel(ch *game.Character, forChannel int32, partyID uint32, targetID uint32, targetName string, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateExpel{
		ForChannel:          forChannel,
		PartyID:             partyID,
		TargetCharacterID:   targetID,
		TargetCharacterName: targetName,
		LeaderCharacterID:   leaderID,
		Members:             members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLeave(ch *game.Character, forChannel int32, partyID uint32, targetID uint32, targetName string, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateLeave{
		ForChannel:          forChannel,
		PartyID:             partyID,
		TargetCharacterID:   targetID,
		TargetCharacterName: targetName,
		LeaderCharacterID:   leaderID,
		Members:             members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateDisband(ch *game.Character, partyID uint32, leaderID uint32) {
	_ = ch.Send(&response.PartyUpdateDisband{
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLeaderChange(ch *game.Character, newLeaderID uint32, byDisconnect bool) {
	_ = ch.Send(&response.PartyUpdateLeaderChange{
		NewLeaderCharacterID: newLeaderID,
		ByDisconnect:         byDisconnect,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLogOnOff(ch *game.Character, forChannel int32, partyID uint32, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateLogOnOff{
		ForChannel:        forChannel,
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
		Members:           members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateSilent(ch *game.Character, forChannel int32, partyID uint32, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateSilent{
		ForChannel:        forChannel,
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
		Members:           members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyStatusMessage(ch *game.Character, code pconst.BuddyStatusCode) {
	ch.Send(&response.BuddyStatus{Code: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyListUpdate(ch *game.Character, action pconst.BuddyListSyncAction, entries []response.BuddyEntry) {
	ch.Send(&response.BuddyListUpdate{
		Action:  action,
		Entries: entries,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyChannelUpdate(ch *game.Character, buddyCharacterID uint32, channel int32) {
	_ = ch.Send(&response.BuddyChannelUpdate{
		CharacterID: buddyCharacterID,
		Channel:     channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyAddRequest(ch *game.Character, fromCharacterID uint32, fromName string) {
	_ = ch.Send(&response.BuddyAddRequest{
		FromCharacterID: fromCharacterID,
		FromName:        fromName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMessage(ch *game.Character, code pconst.GuildResponseCode) {
	_ = ch.Send(&response.GuildMessage{
		Code: code,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnExpGain(ch *game.Character, exp uint32) {
	expPacket := &response.GainExp{
		Gain:  exp,
		White: false,
	}
	ch.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnControlMoveMob(ch *game.Character, mob *game.Mob, moveId uint16, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8) {
	ch.Send(&response.ControlMoveMob{
		OID:          mob.GetOID(),
		MoveId:       moveId,
		EnabledSkill: enabledSkill,
		MP:           mp,
		SkillId:      uint8(skillId),
		SkillLevel:   skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMobHp(ch *game.Character, mob *game.Mob, percentage uint8) {
	ch.Send(&response.ShowMobHp{
		OID:        mob.GetOID(),
		Percentage: percentage,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDojoEnergy(ch *game.Character, energy int) {
	ch.Send(&response.SessionValue{
		Key:   "energy",
		Value: strconv.Itoa(energy),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnOpenNpcShop(ch *game.Character, shopID uint32, shop *wz.Shop) {
	ch.Send(&response.OpenNpcShop{
		ShopID: int32(shopID),
		Shop:   shop,
		Items:  l.gs.GetResources().Items,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) storageTabs(storage *game.Storage, invTypes ...constant.InventoryType) map[constant.InventoryType][]dto.Item {
	tabs := make(map[constant.InventoryType][]dto.Item, len(invTypes))
	for _, invType := range invTypes {
		items := make([]dto.Item, 0, len(storage.Tabs()[invType]))
		for _, item := range storage.Tabs()[invType] {
			items = append(items, item.BuildS2C())
		}
		tabs[invType] = items
	}
	return tabs
}

func (l *CharacterListenerImpl) OnOpenStorage(ch *game.Character, npcID uint32) {
	meso := ch.Account().Storage().Meso()
	ch.Send(&response.Storage{
		Result: pconst.StorageResultOpen,
		NpcID:  npcID,
		Slots:  ch.Account().Storage().Slots(),
		Meso:   &meso,
		Tabs: l.storageTabs(ch.Account().Storage(),
			constant.InventoryTypeEquipment,
			constant.InventoryTypeConsume,
			constant.InventoryTypeInstallation,
			constant.InventoryTypeETC,
			constant.InventoryTypeCash,
		),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageTabChanged(ch *game.Character, result pconst.StorageResult, invType constant.InventoryType) {
	ch.Send(&response.Storage{
		Result: result,
		Slots:  ch.Account().Storage().Slots(),
		Tabs:   l.storageTabs(ch.Account().Storage(), invType),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageArranged(ch *game.Character) {
	ch.Send(&response.Storage{
		Result: pconst.StorageResultArrange,
		Slots:  ch.Account().Storage().Slots(),
		Tabs: l.storageTabs(ch.Account().Storage(),
			constant.InventoryTypeEquipment,
			constant.InventoryTypeConsume,
			constant.InventoryTypeInstallation,
			constant.InventoryTypeETC,
			constant.InventoryTypeCash,
		),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageMesoChanged(ch *game.Character) {
	meso := ch.Account().Storage().Meso()
	ch.Send(&response.Storage{
		Result: pconst.StorageResultMeso,
		Slots:  ch.Account().Storage().Slots(),
		Meso:   &meso,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageError(ch *game.Character, result pconst.StorageResult) {
	ch.Send(&response.StorageError{
		Result: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnOpenDuey(ch *game.Character, fromArrival bool) {
	parcels, expired := ch.Duey().BuildS2C()
	ch.Send(&response.DueyOpen{
		FromArrival: fromArrival,
		Parcels:     parcels,
		Expired:     expired,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyResult(ch *game.Character, result pconst.DueyResult) {
	ch.Send(&response.Duey{
		Result: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyRemoved(ch *game.Character, parcelID uint32, reason uint8) {
	ch.Send(&response.DueyRemoved{
		ParcelID: parcelID,
		Reason:   reason,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyArrival(ch *game.Character, sender string, quick bool, count int) {
	if count > 1 {
		ch.Send(&response.DueyArrivals{
			Quick: quick,
		}, types.SEND_POLICY_ENCRYPT)
		return
	}
	ch.Send(&response.DueyArrival{
		Sender: sender,
		Quick:  quick,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildBulletinThreadList(ch *game.Character, threads []*internal.GuildBulletinBoardThreadEntry, start int, totalCount int, notice *internal.GuildBulletinBoardThreadEntry) {
	entryOf := func(t *internal.GuildBulletinBoardThreadEntry) response.GuildBulletinBoardThreadEntry {
		return response.GuildBulletinBoardThreadEntry{
			LocalThreadID:     t.GetLocalThreadId(),
			PosterCharacterID: t.GetPosterCharacterId(),
			Title:             t.GetTitle(),
			Timestamp:         time.UnixMilli(t.GetTimestampUnixMs()),
			Icon:              t.GetIcon(),
			ReplyCount:        t.GetReplyCount(),
		}
	}
	entries := make([]response.GuildBulletinBoardThreadEntry, 0, len(threads))
	for _, t := range threads {
		entries = append(entries, entryOf(t))
	}
	var noticeEntry *response.GuildBulletinBoardThreadEntry
	if notice != nil {
		n := entryOf(notice)
		noticeEntry = &n
	}
	_ = ch.Send(&response.GuildBulletinBoardThreadList{
		Start:      start,
		TotalCount: totalCount,
		Notice:     noticeEntry,
		Threads:    entries,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildBulletinThread(ch *game.Character, detail *internal.GuildBulletinBoardThreadDetail) {
	replies := make([]response.GuildBulletinBoardReplyEntry, 0, len(detail.GetReplies()))
	for _, r := range detail.GetReplies() {
		replies = append(replies, response.GuildBulletinBoardReplyEntry{
			ReplyID:           r.GetReplyId(),
			PosterCharacterID: r.GetPosterCharacterId(),
			Timestamp:         time.UnixMilli(r.GetTimestampUnixMs()),
			Content:           r.GetContent(),
		})
	}
	_ = ch.Send(&response.GuildBulletinBoardShowThread{
		LocalThreadID:     detail.GetLocalThreadId(),
		PosterCharacterID: detail.GetPosterCharacterId(),
		Timestamp:         time.UnixMilli(detail.GetTimestampUnixMs()),
		Title:             detail.GetTitle(),
		Body:              detail.GetBody(),
		Icon:              detail.GetIcon(),
		Replies:           replies,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShipState(ch *game.Character, state uint16) {
	ch.Send(&response.ShipState{State: state}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShipBalrog(ch *game.Character) {
	ch.Send(&response.ShipSpecialEffect{Effect: response.ShipSpecialBalrog}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSitOnChair(ch *game.Character, itemID uint32) {
	ch.Broadcast(&response.ShowChair{
		CharacterID: ch.GetID(),
		ItemID:      itemID,
	}, nil)
}

func (l *CharacterListenerImpl) OnSitOnMapSeat(ch *game.Character, seatID int16) {
	ch.Send(&response.CancelChair{
		ChairID: seatID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStandUp(ch *game.Character) {
	ch.Send(&response.CancelChair{
		ChairID: -1,
	}, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(&response.ShowChair{
		CharacterID: ch.GetID(),
		ItemID:      0,
	}, nil)
}

func (l *CharacterListenerImpl) OnUnlockAction(ch *game.Character) {
	ch.Send(&response.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInspect(ch *game.Character, profile *response.CharacterProfile) {
	ch.Send(profile, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCardRegistered(ch *game.Character, cardID uint32, count uint32) {
	ch.Send(&response.MonsterBookSetCard{Success: true, CardID: cardID, Count: count}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCardFull(ch *game.Character) {
	ch.Send(&response.MonsterBookSetCard{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCover(ch *game.Character, cardID uint32) {
	ch.Send(&response.MonsterBookSetCover{CardID: cardID}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTeleportStones(ch *game.Character, vip bool) {
	ch.Send(&response.TeleportStoneResult{Result: pconst.TeleportStoneResultList, VIP: vip, Maps: ch.TeleportStones().Slots(vip)}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTeleportStoneFailed(ch *game.Character, vip bool, result pconst.TeleportStoneResult) {
	ch.Send(&response.TeleportStoneResult{Result: result, VIP: vip}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnScriptError(ch *game.Character, script string, err error) {
	if ch.GetRole() != constant.RoleAdmin {
		return
	}
	l.OnMessage(ch, constant.MsgPinkText, fmt.Sprintf("스크립트 오류 %s: %v", script, err))
}

func (l *CharacterListenerImpl) OnQuestStarted(ch *game.Character, qp *game.Quest, npcID uint32) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
	if npcID != 0 {
		ch.Send(&response.UpdateQuestNPC{
			Progress:    8,
			QuestID:     status.QuestID,
			NPCID:       npcID,
			NextQuestID: 0,
		}, types.SEND_POLICY_ENCRYPT)
	}
	if limit := qp.TimeLimit(); limit > 0 {
		l.OnClock(ch, int32(limit))
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestCompleted(ch *game.Character, qp *game.Quest, npcID uint32, nextQuestID uint32) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
	if npcID != 0 {
		ch.Send(&response.UpdateQuestNPC{
			Progress:    8,
			QuestID:     status.QuestID,
			NPCID:       npcID,
			NextQuestID: nextQuestID,
		}, types.SEND_POLICY_ENCRYPT)
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestForfeited(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnShowQuestCompletion(ch *game.Character, questID uint32) {
	if questID == 0 || questID > 0xFFFF {
		return
	}
	ch.Send(&response.ShowQuestCompletion{
		QuestID: uint16(questID),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPlaySound(ch *game.Character, sound string, broadcast bool) {
	if sound == "" {
		return
	}
	pkt := &response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeSound,
		Env:  sound,
	}
	if broadcast {
		if m := ch.GetMap(); m != nil {
			m.Broadcast(pkt, nil)
			return
		}
	}
	ch.Send(pkt, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPlayPortalSound(ch *game.Character) {
	ch.Send(&response.ShowSelfSkillEffect{
		Type:       pconst.SkillEffectTypePortal,
		SkillID:    0,
		SkillLevel: 1,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnQuestProgress(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	status.Status = response.QuestWireStatusStarted
	ch.Send(&response.UpdateQuest{QuestStatus: status}, types.SEND_POLICY_ENCRYPT)

	if qp.CanComplete(ch, game.QuestPhaseOpts{}) == nil {
		ch.Send(&response.ShowQuestCompletion{
			QuestID: status.QuestID,
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnQuestRecordExChanged(ch *game.Character, qp *game.Quest) {
	status, ok := qp.BuildS2C()
	if !ok {
		return
	}
	ch.Send(&response.UpdateQuestRecordEx{
		QuestID: status.QuestID,
		Data:    qp.RecordExWire(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemGainFailed(ch *game.Character, mode constant.ItemGainFailedType) {
	ch.Send(&response.ItemGainFailed{
		Mode: mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotUpdated(ch *game.Character, inventoryType constant.InventoryType, slot int16, item game.Item) {
	var itemDTO dto.Item
	if item != nil {
		itemDTO = item.BuildS2C()
	}
	ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotAdded(ch *game.Character, inventoryType constant.InventoryType, slot int16, item game.Item) {
	var itemDTO dto.Item
	if item != nil {
		itemDTO = item.BuildS2C()
	}
	ch.Send(&response.AddInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
		FromDrop:      true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowItemGain(ch *game.Character, itemId uint32, count uint32, mode constant.ShowItemGainType) {
	ch.Send(&response.ShowItemGain{
		ItemId: itemId,
		Count:  count,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMesoGain(ch *game.Character, count int32, mode constant.ShowMesoGainType) {
	ch.Send(&response.ShowMesoGain{
		Count: count,
		Mode:  mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateStats(ch *game.Character, stats map[constant.Stat]int32, unlock bool) {
	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: unlock,
	}, types.SEND_POLICY_ENCRYPT)
	if stats != nil {
		_, hasHP := stats[constant.StatHP]
		_, hasMaxHP := stats[constant.StatMaxHP]
		if m := ch.GetMap(); m != nil && (hasHP || hasMaxHP) {
			m.BroadcastPartyHP(ch)
		}
	}
}

func (l *CharacterListenerImpl) OnShowSelfSkillEffect(ch *game.Character, effectType pconst.SkillEffectType, skillID uint32, skillLevel uint8, additional *uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfSkillEffect{
		Type:       effectType,
		SkillID:    skillID,
		SkillLevel: skillLevel,
		Additional: additional,
	}, types.SEND_POLICY_ENCRYPT)

	l.OnShowSkillEffect(ch, effectType, skillID, skillLevel, additional)
}

func (l *CharacterListenerImpl) OnShowSkillEffect(ch *game.Character, effectType pconst.SkillEffectType, skillID uint32, skillLevel uint8, additional *uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowSkillEffect{
		CharacterID: ch.GetID(),
		Type:        effectType,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
		Additional:  additional,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfEffect(ch *game.Character, effectType response.EffectType) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfEffect{
		Type: effectType,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowEffect(ch, effectType)
}

func (l *CharacterListenerImpl) OnShowEffect(ch *game.Character, effectType response.EffectType) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowEffect{
		CharacterID: ch.GetID(),
		Type:        effectType,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfDragonBloodEffect(ch *game.Character, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfDragonBloodEffect{
		SkillID:    skillID,
		SkillLevel: skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowDragonBloodEffect(ch, skillID, skillLevel)
}

func (l *CharacterListenerImpl) OnShowDragonBloodEffect(ch *game.Character, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowDragonBloodEffect{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfHPHealedEffect(ch *game.Character, amount int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfHPHealedEffect{
		Amount: amount,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowHPHealedEffect(ch, amount)
}

func (l *CharacterListenerImpl) OnShowHPHealedEffect(ch *game.Character, amount int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowHPHealedEffect{
		CharacterID: ch.GetID(),
		Amount:      amount,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfRewardItemAnimation(ch *game.Character, itemID uint32, effect string) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfRewardItemAnimation{
		ItemID: itemID,
		Effect: effect,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowRewardItemAnimation(ch, itemID, effect)
}

func (l *CharacterListenerImpl) OnShowRewardItemAnimation(ch *game.Character, itemID uint32, effect string) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowRewardItemAnimation{
		CharacterID: ch.GetID(),
		ItemID:      itemID,
		Effect:      effect,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfItemMakerSuccessEffect(ch *game.Character) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfItemMakerSuccessEffect{}, types.SEND_POLICY_ENCRYPT)
	l.OnShowItemMakerSuccessEffect(ch)
}

func (l *CharacterListenerImpl) OnShowItemMakerSuccessEffect(ch *game.Character) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowItemMakerSuccessEffect{
		CharacterID: ch.GetID(),
	}, nil)
}

func (l *CharacterListenerImpl) OnItemMakerDestroyed(ch *game.Character) {
	ch.Send(&response.ShowSelfItemMakerSuccessEffect{}, types.SEND_POLICY_ENCRYPT)
	l.OnMessage(ch, constant.MsgPopup, "아이템이 파괴되었습니다.")
}

func (l *CharacterListenerImpl) OnUseSkillBook(ch *game.Character, skillID uint32, masterLevel int, canUse bool, success bool) {
	ch.Broadcast(&response.UseSkillBook{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		MasterLevel: uint32(masterLevel),
		CanUse:      canUse,
		Success:     success,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnShowSelfCraftingEffect(ch *game.Character, effect string, time int32, mode int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfCraftingEffect{
		Effect: effect,
		Time:   time,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowCraftingEffect(ch, effect, time, mode)
}

func (l *CharacterListenerImpl) OnShowCraftingEffect(ch *game.Character, effect string, time int32, mode int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowCraftingEffect{
		CharacterID: ch.GetID(),
		Effect:      effect,
		Time:        time,
		Mode:        mode,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfDiceEffect(ch *game.Character, effectID int32, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfDiceEffect{
		EffectID:   effectID,
		SkillID:    skillID,
		SkillLevel: skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowDiceEffect(ch, effectID, skillID, skillLevel)
}

func (l *CharacterListenerImpl) OnShowDiceEffect(ch *game.Character, effectID int32, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowDiceEffect{
		CharacterID: ch.GetID(),
		EffectID:    effectID,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowScrollEffect(ch *game.Character, success bool, destroyedByCurse bool) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowScrollEffect{
		CharacterID:      ch.GetID(),
		Success:          success,
		DestroyedByCurse: destroyedByCurse,
	}, &game.ObjectBroadcastOption{
		WithMe: true,
	})
}

func (l *CharacterListenerImpl) OnScrolledItem(ch *game.Character, scrollInventoryType constant.InventoryType, scrollSlot int16, scrollCount uint16, upgradedSlot int16, destroyed bool, potential bool, upgradedItem game.Item) {
	var upgradedItemDTO dto.Item
	if upgradedItem != nil {
		upgradedItemDTO = upgradedItem.BuildS2C()
	}
	ch.Send(&response.ScrolledItem{
		ScrollInventoryType: scrollInventoryType,
		ScrollSlot:          scrollSlot,
		ScrollCount:         scrollCount,
		UpgradedSlot:        upgradedSlot,
		Destroyed:           destroyed,
		Potential:           potential,
		UpgradedItem:        upgradedItemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMobMoved(ch *game.Character, mob *game.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}

	movePacket := &response.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mob.GetOID(),
		StartPoint:  startPoint,
		Movements:   movements,
	}

	controllerID := uint32(0)
	if ch != nil {
		controllerID = ch.GetID()
	}
	for _, obj := range mob.Nears(constant.ObjectTypeCharacter, &game.SearchOption{
		IncludeHidden: true,
	}) {
		peer, ok := obj.(*game.Character)
		if !ok || peer == nil || peer.GetID() == controllerID {
			continue
		}
		_ = peer.Send(movePacket, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnPlayerMove(ch *game.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	if ch.GetMap() == nil {
		return
	}

	characterDTO := ch.BuildS2C()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	ch.Broadcast(movePacket, nil)
}

func (l *CharacterListenerImpl) OnFieldRelocate(ch *game.Character, spawnPoint uint8) {
	ch.Send(&response.FieldRelocate{Portal: spawnPoint}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAttack(ch *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.Attack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnEnergyAttack(ch *game.Character, attackPayload dto.AttackPayload) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.EnergyAttack{
		Attack: response.Attack{
			AttackInfo:  attackPayload.ToAttackInfo(),
			CharacterId: ch.GetID(),
		},
	}, nil)
}

func (l *CharacterListenerImpl) OnRangedAttack(ch *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.RangedAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
		CashBullet:  0,
	}, nil)
}

func (l *CharacterListenerImpl) OnMagicAttack(ch *game.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.MagicAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnEndSortInventory(ch *game.Character, inventoryType constant.InventoryType) {
	ch.Send(&response.EndSortInventory{
		InventoryType: inventoryType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSwapInventorySlot(ch *game.Character, inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8) {
	ch.Send(&response.SwapInventorySlot{
		InventoryType:   inventoryType,
		Source:          source,
		Dest:            dest,
		EquipmentAction: response.EquipmentActionType(equipmentAction),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnRemoveInventorySlot(ch *game.Character, inventoryType constant.InventoryType, slot int16) {
	ch.Send(&response.RemoveInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemExpired(ch *game.Character, items []wz.Item) {
	var expired []uint32
	for _, model := range items {
		if model.GetTraits()&constant.ItemTraitCash != 0 {
			ch.Send(&response.CashItemExpired{ItemID: model.GetID()}, types.SEND_POLICY_ENCRYPT)
		} else {
			expired = append(expired, model.GetID())
		}
	}
	if len(expired) == 0 {
		return
	}
	ch.Send(&response.ItemExpired{ItemIDs: expired}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFullMergeInventorySlot(ch *game.Character, inventoryType constant.InventoryType, source int16, dest int16, count uint16) {
	ch.Send(&response.FullMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		Count:         count,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartialMergeInventorySlot(ch *game.Character, inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16) {
	ch.Send(&response.PartialMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		SourceCount:   sourceCount,
		DestCount:     destCount,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateCharacterLook(ch *game.Character) {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	lookPacket := &response.UpdateCharacterLook{
		Character:      ch.BuildS2C(),
		CrushRing:      ch.Inventory().WornRing(ch.Inventory().Rings().Left()),
		FriendshipRing: ch.Inventory().WornRing(ch.Inventory().Rings().Mid()),
		MarriageRing:   ch.Wedding().BuildRingS2C(),
	}

	ch.Broadcast(lookPacket, nil)
	if ch.Messenger().RoomID() != 0 {
		l.gs.systems.Get[*game.DispatchSystem]().Visit(ch.GetID(), func(ctx actor.Context, current *game.Character) {
			if current.Messenger().RoomID() == 0 {
				return
			}
			w := stream.NewStreamWriter(stream.LittleEndian)
			current.BuildS2C().SerializeLook(w)
			req := &internal.MessengerRequest{
				WorldId:     l.gs.config.WorldId,
				CharacterId: current.GetID(),
				ChannelId:   l.gs.config.ChannelId,
				Action:      internal.MessengerAction_MESSENGER_ACTION_AVATAR,
				Look:        w.Bytes(),
			}
			async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
				func(c context.Context) (*internal.MessengerReply, error) {
					return l.gs.systems.Get[*rpc.InternalConnection]().Client.Messenger(c, req)
				},
				nil,
			).OnError(func(error) {})
		})
	}
}

func (l *CharacterListenerImpl) OnNpcAction(ch *game.Character, bytes []byte) {
	ch.Send(&response.NpcAction{
		Bytes: bytes,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnClassChange(ch *game.Character, oldClass constant.ClassType, newClass constant.ClassType, skillPoint uint16) {
	stats := map[constant.Stat]int32{
		constant.StatClass:       int32(newClass),
		constant.StatAvailableSP: int32(skillPoint),
	}

	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
	l.gs.systems.Get[*game.PartySystem]().UpdateMemberAsync(nil, ch)
	l.gs.systems.Get[*game.AllianceSystem]().NotifyMemberFieldsChanged(ch)
}

func (l *CharacterListenerImpl) OnPartyMemberFieldsChanged(ch *game.Character) {
	l.gs.systems.Get[*game.PartySystem]().UpdateMemberAsync(nil, ch)
	l.gs.systems.Get[*game.AllianceSystem]().NotifyMemberFieldsChanged(ch)
}

func (l *CharacterListenerImpl) OnPartyMemberHP(recipient *game.Character, member *game.Character) {
	_ = recipient.Send(&response.UpdatePartyMemberHP{
		CharacterID: member.GetID(),
		CurrentHP:   int32(member.GetHp()),
		MaxHP:       int32(member.GetMaxHp()),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuffAdded(ch *game.Character, buffID int32, remainingDuration time.Duration, values map[constant.BuffFlag]int32) {
	if len(values) == 0 {
		return
	}
	dtoBuffs := make([]dto.BuffEntry, 0, len(values))
	for flag, value := range values {
		dtoBuff := dto.BuffEntry{Buff: flag, Value: value}
		dtoBuffs = append(dtoBuffs, dtoBuff)
	}

	var selfPacket types.Packet
	var remotePacket types.Packet
	if mountID, ok := values[constant.BuffFlagMonsterRiding]; ok {
		selfPacket = &response.UpdateSelfRidding{
			BuffID:  buffID,
			MountID: mountID,
			Buffs:   dtoBuffs,
		}
		remotePacket = &response.UpdateRidding{
			CharacterID: int32(ch.GetID()),
			MountID:     mountID,
			SkillID:     buffID,
			Buffs:       dtoBuffs,
		}
	} else {
		selfPacket = &response.UpdateSelfBuff{
			BuffID:   buffID,
			Duration: remainingDuration,
			Buffs:    dtoBuffs,
		}
		remotePacket = &response.UpdateBuff{
			CharacterID: int32(ch.GetID()),
			BuffID:      buffID,
			Duration:    remainingDuration,
			Buffs:       dtoBuffs,
		}
	}

	ch.Send(selfPacket, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(remotePacket, nil)
}

func (l *CharacterListenerImpl) OnBuffRemoved(ch *game.Character, flags []constant.BuffFlag) {
	ch.Send(&response.CancelSelfBuff{Buffs: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.CancelBuff{
		CharacterID: int32(ch.GetID()),
		Buffs:       flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffAdded(ch *game.Character, debuff constant.DebuffFlag, x int16, skillID uint16, skillLevel uint16, durationMs int32) {
	ch.Send(&response.GiveSelfDebuff{
		Debuff:     debuff,
		X:          x,
		SkillID:    skillID,
		SkillLevel: skillLevel,
		DurationMs: durationMs,
	}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.GiveDebuff{
		CharacterID: int32(ch.GetID()),
		Debuff:      debuff,
		X:           x,
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffRemoved(ch *game.Character, flags []constant.DebuffFlag) {
	ch.Send(&response.RemoveSelfDebuff{Diseases: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.RemoveDebuff{
		CharacterID: int32(ch.GetID()),
		Diseases:    flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnUpdateSkill(ch *game.Character, skillID uint32, level int32, masterLevel int32) {
	ch.Send(&response.UpdateSkills{
		SkillID:     skillID,
		Level:       level,
		MasterLevel: masterLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSkillCooldown(ch *game.Character, skillID uint32, remainingSec uint16) {
	ch.Send(&response.SkillCooldown{
		SkillID:      skillID,
		RemainingSec: uint32(remainingSec),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnHiddenChanged(ch *game.Character, hidden bool) {
	ch.Send(&response.SuperHide{Hidden: hidden}, types.SEND_POLICY_ENCRYPT)

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	if hidden {
		ch.Broadcast(&response.LeavePlayer{ID: ch.GetID()}, &game.ObjectBroadcastOption{
			RecipientsRoleBelowPivot: true,
		})
	} else {
		for _, obj := range mapInstance.GetObjectsNear(ch.GetPosition(), constant.ObjectTypeCharacter, &game.SearchOption{IncludeHidden: true}) {
			if viewer, ok := obj.(*game.Character); ok {
				ch.SendSpawnSyncToViewer(viewer)
			}
		}
	}
	if controllerTable := mapInstance.GetControllerTable(); controllerTable != nil {
		controllerTable.Update(ch)
	}
}

func (l *CharacterListenerImpl) OnPetSpawn(ch *game.Character, sn uint64) {
	ch.Broadcast(&response.SpawnPet{
		CharacterID: ch.GetID(),
		Pet:         ch.Pets().Active().BuildS2C(),
	}, &game.ObjectBroadcastOption{WithMe: true})
	ch.Pets().ShowExceptions()
	ch.Send(&response.UpdateStats{Pet: &sn, UnlockAction: true}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetRemove(ch *game.Character, reason constant.PetRemoveReason) {
	ch.Broadcast(&response.SpawnPet{
		CharacterID: ch.GetID(),
		Reason:      reason,
	}, &game.ObjectBroadcastOption{WithMe: true})
	sn := uint64(0)
	ch.Send(&response.UpdateStats{Pet: &sn, UnlockAction: true}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetMove(ch *game.Character, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	ch.Broadcast(&response.MovePet{
		CharacterID: ch.GetID(),
		StartPoint:  startPoint,
		Fragments:   movements,
	}, nil)
}

func (l *CharacterListenerImpl) OnPetChat(ch *game.Character, typ uint8, action uint8, text string) {
	ch.Broadcast(&response.PetChat{
		CharacterID: ch.GetID(),
		Type:        typ,
		Action:      action,
		Text:        text,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetCommand(ch *game.Character, index uint8, success bool) {
	ch.Broadcast(&response.PetCommand{CharacterID: ch.GetID(), Index: index, Success: success}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetFood(ch *game.Character, success bool) {
	ch.Broadcast(&response.PetCommand{CharacterID: ch.GetID(), Food: true, Success: success}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnLieDetectorStarted(ch *game.Character, target string) {
	ch.Send(&response.LieDetectorStarted{Target: target}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnLieDetectorQuestion(ch *game.Character, image []byte, skill bool) {
	ch.Send(&response.LieDetectorQuestion{Image: image, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnLieDetectorResult(ch *game.Character, result pconst.LieDetectorResult, skill bool) {
	ch.Send(&response.LieDetectorResult{Result: result, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetUpdated(ch *game.Character, pet *game.Pet) {
	slot, ok := ch.Inventory().FindSlot(constant.InventoryTypeCash, pet)
	if !ok {
		return
	}
	l.OnInventorySlotUpdated(ch, constant.InventoryTypeCash, slot, pet)
}

func (l *CharacterListenerImpl) OnMountUpdated(ch *game.Character, level uint32, exp uint32, fatigue uint32, levelUp bool) {
	ch.Broadcast(&response.UpdateMount{
		CharacterID: ch.GetID(),
		Level:       level,
		Exp:         exp,
		Fatigue:     fatigue,
		LevelUp:     levelUp,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetExceptions(ch *game.Character, sn uint64, itemIDs []uint32) {
	ch.Send(&response.PetExceptions{
		CharacterID: ch.GetID(),
		SN:          sn,
		ItemIDs:     itemIDs,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetNameChanged(ch *game.Character, name string) {
	ch.Broadcast(&response.PetNameChanged{
		CharacterID: ch.GetID(),
		Name:        name,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetSkillChanged(ch *game.Character, pet *game.Pet, sn uint64, skill constant.PetSkill, add bool) {
	l.OnPetUpdated(ch, pet)
	ch.Send(&response.PetSkillChanged{SN: sn, Add: add, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetLevelUp(ch *game.Character) {
	ch.Send(&response.ShowSelfEffect{Type: response.EffectTypePet}, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(&response.ShowEffect{CharacterID: ch.GetID(), Type: response.EffectTypePet}, nil)
}

func (l *CharacterListenerImpl) OnSummonSpawn(ch *game.Character, summon *game.Summon) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SpawnSummon{
		OwnerID:      ch.GetID(),
		OID:          summon.GetOID(),
		SkillID:      summon.SkillID(),
		SkillLevel:   summon.SkillLevel(),
		Position:     summon.GetPosition(),
		MovementType: summon.MovementType(),
		SummonType:   summon.SummonType(),
		Animated:     true,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonRemove(ch *game.Character, summon *game.Summon, animated bool) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.RemoveSummon{
		OwnerID:  ch.GetID(),
		OID:      summon.GetOID(),
		Animated: animated,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonMove(ch *game.Character, summon *game.Summon, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.MoveSummon{
		CharacterID: ch.GetID(),
		OID:         summon.GetOID(),
		StartPoint:  startPoint,
		Fragments:   movements,
	}
	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnSummonAttack(ch *game.Character, summon *game.Summon, animation uint8, targets []game.SummonAttackTarget) {
	if summon.GetMap() == nil {
		return
	}
	respTargets := make([]response.SummonAttackTarget, len(targets))
	for i, t := range targets {
		respTargets[i] = response.SummonAttackTarget{
			OID:    t.OID,
			Damage: t.Damage,
		}
	}
	packet := &response.SummonAttack{
		CharacterID:   ch.GetID(),
		SummonSkillID: uint32(summon.SkillID()),
		Animation:     animation,
		Targets:       respTargets,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonSkill(ch *game.Character, summon *game.Summon, newStance uint8) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SummonSkill{
		CharacterID: ch.GetID(),
		SummonOID:   summon.GetOID(),
		NewStance:   newStance,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonDamaged(ch *game.Character, summon *game.Summon, unknown uint8, damage uint32, monsterIdFrom uint32) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.DamageSummon{
		CharacterID:   ch.GetID(),
		SummonSkillID: uint32(summon.SkillID()),
		Unknown:       unknown,
		Damage:        damage,
		MonsterIDFrom: monsterIdFrom,
	}
	ch.Broadcast(packet, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnCarnivalStart(ch *game.Character, team constant.CarnivalTeam, persAvail, persTotal, friendAvail, friendTotal, enemyAvail, enemyTotal int) {
	ch.Send(&response.CarnivalStart{
		Team:                team,
		PersonalAvailableCP: uint16(persAvail),
		PersonalTotalCP:     uint16(persTotal),
		TeamAvailableCP:     uint16(friendAvail),
		TeamTotalCP:         uint16(friendTotal),
		EnemyAvailableCP:    uint16(enemyAvail),
		EnemyTotalCP:        uint16(enemyTotal),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalObtainedCP(ch *game.Character, avail, total int) {
	ch.Send(&response.CarnivalObtainedCP{
		Avail: uint16(avail),
		Total: uint16(total),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalPartyCP(ch *game.Character, team constant.CarnivalTeam, avail, total int) {
	ch.Send(&response.CarnivalPartyCP{
		Team:  team,
		Avail: uint16(avail),
		Total: uint16(total),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalSummon(ch *game.Character, tab constant.CarnivalTab, num uint8, name string) {
	ch.Broadcast(&response.CarnivalSummon{
		Tab:  uint8(tab),
		Num:  num,
		Name: name,
	}, &game.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnCarnivalDied(ch *game.Character, team constant.CarnivalTeam, name string, lostCP uint8) {
	ch.Send(&response.CarnivalDied{
		Team:   team,
		Name:   name,
		LostCP: lostCP,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalResult(ch *game.Character, winner bool) {
	message := "아쉽게도 패배하였습니다. 잠시 후 자동으로 보상 맵으로 나가집니다."
	effect := "quest/carnival/lose"
	sound := "MobCarnival/Lose"
	if winner {
		message = "승리하였습니다! 잠시 후 자동으로 보상 맵으로 나가집니다."
		effect = "quest/carnival/win"
		sound = "MobCarnival/Win"
	}
	l.OnMessage(ch, constant.MsgLightBlueText, message)
	ch.Send(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeMapEffect,
		Env:  effect,
	}, types.SEND_POLICY_ENCRYPT)
	ch.Send(&response.EnvironmentChange{
		Mode: response.EnvironmentChangeModeSound,
		Env:  sound,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnConfirmShopTransaction(ch *game.Character, code uint8) {
	_ = ch.Send(&response.ConfirmShopTransaction{
		Code: code,
	}, types.SEND_POLICY_ENCRYPT)
}

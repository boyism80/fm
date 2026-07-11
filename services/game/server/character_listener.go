package server

import (
	"time"

	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type CharacterListenerImpl struct {
	gs *GameServer
}

func (l *CharacterListenerImpl) OnDialog(ch *entity.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.Dialog{
		NPC:  npc,
		Type: constant.DialogTypeDefault,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogYesNo(ch *entity.Character, npc uint32, message string, prev bool, next bool) {
	dialogPacket := &response.DialogYesNo{
		NPC:  npc,
		Text: message,
		Prev: prev,
		Next: next,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogAccept(ch *entity.Character, npc uint32, message string, enableEscape bool) {
	dialogPacket := &response.DialogAccept{
		NPC:          npc,
		Text:         message,
		EnableEscape: enableEscape,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogList(ch *entity.Character, npc uint32, message string, selections []string) {
	dialogPacket := &response.DialogList{
		NPC:        npc,
		Text:       message,
		Selections: selections,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDialogInput(ch *entity.Character, npc uint32, message string) {
	dialogPacket := &response.DialogInput{
		NPC:  npc,
		Text: message,
	}
	ch.Send(dialogPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnChat(ch *entity.Character, message string, highlight bool, dontRecordHistory bool) {
	chatPacket := &response.NormalChat{
		CharacterId:       ch.GetID(),
		Message:           message,
		Highlight:         highlight,
		DontRecordHistory: dontRecordHistory,
	}

	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(chatPacket, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnMesoChanged(ch *entity.Character, meso int32) {
	ch.Send(&response.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessage(ch *entity.Character, messageType constant.ServerMessageType, message string) {
	noticePacket := &response.Notice{
		Message: message,
		Type:    messageType,
	}
	ch.Send(noticePacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyCreated(ch *entity.Character, partyID uint32) {
	if ch == nil {
		return
	}
	ch.Send(&response.PartyCreated{
		PartyID: partyID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowGuildInfo(ch *entity.Character) {
	if l == nil || l.gs == nil || ch == nil {
		return
	}
	guildID, ok := ch.GetGuildID()
	if !ok {
		return
	}
	ent := l.gs.guild.Get(guildID)
	if ent == nil {
		return
	}
	info := entity.GuildToDTO(ent)
	if info == nil {
		return
	}
	ch.Send(&response.GuildShowInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowAllianceInfo(ch *entity.Character) {
	if l == nil || l.gs == nil || ch == nil {
		return
	}
	allianceID := uint32(0)
	if guildID, ok := ch.GetGuildID(); ok {
		if g := l.gs.guild.Get(guildID); g != nil {
			if id, inAlliance := g.GetAllianceID(); inAlliance {
				allianceID = id
			}
		}
	}
	info, guilds := l.gs.alliance.ShowData(allianceID)
	_ = ch.Send(&response.AllianceShowInfo{Info: info}, types.SEND_POLICY_ENCRYPT)
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceShowGuilds{Guilds: guilds}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceUpdateInfo{Info: info}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceUpdateInfo(ch *entity.Character) {
	if l == nil || l.gs == nil || ch == nil {
		return
	}
	allianceID := uint32(0)
	if guildID, ok := ch.GetGuildID(); ok {
		if g := l.gs.guild.Get(guildID); g != nil {
			if id, inAlliance := g.GetAllianceID(); inAlliance {
				allianceID = id
			}
		}
	}
	if allianceID == 0 {
		return
	}
	info, _ := l.gs.alliance.ShowData(allianceID)
	if info == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateInfo{Info: info}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBroadcastGuildAppearance(ch *entity.Character) {
	if l == nil || ch == nil {
		return
	}
	guildID, ok := ch.GetGuildID()
	if !ok || l.gs == nil {
		return
	}
	ent := l.gs.guild.Get(guildID)
	if ent == nil {
		return
	}
	characterID := ch.GetID()
	broadcastOpt := &entity.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   ent.Name,
	}, broadcastOpt)
	iconPkt := &response.LoadGuildIcon{
		CharacterID: characterID,
	}
	if ent.Logo != nil {
		iconPkt.HasGuildIcon = true
		iconPkt.LogoBG = uint16(ent.Logo.LogoBG)
		iconPkt.LogoBGColor = uint8(ent.Logo.LogoBGColor)
		iconPkt.Logo = uint16(ent.Logo.Logo)
		iconPkt.LogoColor = uint8(ent.Logo.LogoColor)
	}
	ch.Broadcast(iconPkt, broadcastOpt)
}

func (l *CharacterListenerImpl) OnPartyInvite(ch *entity.Character, partyID uint32, inviterName string, partySearch bool) {
	if ch == nil {
		return
	}
	ch.Send(&response.PartyInvite{
		PartyID:     partyID,
		InviterName: inviterName,
		PartySearch: partySearch,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildInvite(ch *entity.Character, guildID uint32, inviterName string) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildInvite{
		GuildID:     guildID,
		InviterName: inviterName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceInvite(ch *entity.Character, inviterGuildID uint32, inviterName string, allianceName string) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.AllianceInvite{
		InviterGuildID: inviterGuildID,
		InviterName:    inviterName,
		AllianceName:   allianceName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceCreate(ch *entity.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, membershipGuilds []dto.AllianceMembershipChangeGuild) {
	if ch == nil || info == nil {
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

func (l *CharacterListenerImpl) OnAllianceDisband(ch *entity.Character, allianceID uint32) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.AllianceDisband{
		AllianceID: allianceID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceInfoBroadcast(ch *entity.Character, info *dto.AllianceInfo) {
	if ch == nil || info == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceNoticeChanged(ch *entity.Character, info *dto.AllianceInfo) {
	l.OnAllianceInfoBroadcast(ch, info)
}

func (l *CharacterListenerImpl) OnAllianceLeaderChanged(ch *entity.Character, allianceID uint32, oldLeaderID uint32, newLeaderID uint32, info *dto.AllianceInfo, guilds []*dto.GuildInfo) {
	if ch == nil || info == nil {
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

func (l *CharacterListenerImpl) OnAllianceMemberRankChanged(ch *entity.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo) {
	if ch == nil || info == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateInfo{
		Info: info,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.AllianceShowGuilds{
		Guilds: guilds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceGuildAdded(ch *entity.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, newGuildID uint32, addedGuild *dto.GuildInfo, members []dto.AllianceGuildMemberRank, joining bool, membershipGuild *dto.AllianceMembershipChangeGuild) {
	if ch == nil || info == nil || addedGuild == nil {
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

func (l *CharacterListenerImpl) OnAllianceGuildLeft(ch *entity.Character, info *dto.AllianceInfo, removedGuildID uint32, removedGuild *dto.GuildInfo, removedMembers []dto.AllianceGuildMemberRank, expelled bool, leaving bool) {
	if ch == nil || info == nil || removedGuild == nil {
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

func (l *CharacterListenerImpl) OnGuildNewMember(ch *entity.Character, guildID uint32, member dto.GuildMemberStatus) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildNewMember{
		GuildID: guildID,
		Member:  member,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildLeaveSelf(ch *entity.Character) {
	if l == nil || ch == nil {
		return
	}
	ch.SetGuildID(nil)
	_ = ch.Send(&response.GuildShowInfo{
		Info: nil,
	}, types.SEND_POLICY_ENCRYPT)
	characterID := ch.GetID()
	broadcastOpt := &entity.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   "",
	}, broadcastOpt)
	ch.Broadcast(&response.LoadGuildIcon{
		CharacterID: characterID,
	}, broadcastOpt)
}

func (l *CharacterListenerImpl) OnGuildExpelledSelf(ch *entity.Character, guildID uint32) {
	if l == nil || ch == nil {
		return
	}
	ch.SetGuildID(nil)
	_ = ch.Send(&response.GuildMemberLeft{
		GuildID:     guildID,
		CharacterID: ch.GetID(),
		Name:        ch.GetName(),
		WasExpelled: true,
	}, types.SEND_POLICY_ENCRYPT)
	characterID := ch.GetID()
	broadcastOpt := &entity.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   "",
	}, broadcastOpt)
	ch.Broadcast(&response.LoadGuildIcon{
		CharacterID: characterID,
	}, broadcastOpt)
}

func (l *CharacterListenerImpl) OnGuildMemberLeft(ch *entity.Character, guildID uint32, targetID uint32, targetName string, wasExpelled bool) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildMemberLeft{
		GuildID:     guildID,
		CharacterID: targetID,
		Name:        targetName,
		WasExpelled: wasExpelled,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildRankTitleChange(ch *entity.Character, guildID uint32, rankTitles [5]string) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildRankTitleChange{
		GuildID:    guildID,
		RankTitles: rankTitles,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberRankChange(ch *entity.Character, guildID uint32, targetID uint32, guildRank uint8) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildChangeRank{
		GuildID:     guildID,
		CharacterID: targetID,
		GuildRank:   guildRank,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildEmblemChange(ch *entity.Character, guildID uint32, logoBG uint16, logoBGColor uint8, logo uint16, logoColor uint8) {
	if l == nil || ch == nil {
		return
	}
	_ = ch.Send(&response.GuildEmblemChange{
		GuildID:     guildID,
		LogoBG:      logoBG,
		LogoBGColor: logoBGColor,
		Logo:        logo,
		LogoColor:   logoColor,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnBroadcastGuildAppearance(ch)
}

func (l *CharacterListenerImpl) OnGuildNoticeChange(ch *entity.Character, guildID uint32, notice string) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildNotice{
		GuildID: guildID,
		Notice:  notice,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildCapacityChange(ch *entity.Character, guildID uint32, capacity uint8) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildCapacityChange{
		GuildID:  guildID,
		Capacity: capacity,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberOnlineChange(ch *entity.Character, guildID uint32, subjectCharacterID uint32, online bool) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildMemberOnline{
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Online:      online,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberFieldsChange(ch *entity.Character, guildID uint32, subjectCharacterID uint32, level uint32, classID uint32) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildMemberLevelClassUpdate{
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Level:       level,
		ClassID:     classID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberOnlineChange(ch *entity.Character, allianceID uint32, guildID uint32, subjectCharacterID uint32, online bool) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.AllianceMemberOnline{
		AllianceID:  allianceID,
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Online:      online,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberFieldsChange(ch *entity.Character, allianceID uint32, guildID uint32, subjectCharacterID uint32, level uint32, classID uint32) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.AllianceUpdateMember{
		AllianceID:  allianceID,
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Level:       level,
		ClassID:     classID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildDisbandSelf(ch *entity.Character, guildID uint32) {
	if l == nil || ch == nil {
		return
	}
	if gid, ok := ch.GetGuildID(); ok && gid == guildID {
		_ = ch.Send(&response.GuildDisband{
			GuildID: guildID,
		}, types.SEND_POLICY_ENCRYPT)
	}
	ch.SetGuildID(nil)
	_ = ch.Send(&response.GuildShowInfo{
		Info: nil,
	}, types.SEND_POLICY_ENCRYPT)
	characterID := ch.GetID()
	broadcastOpt := &entity.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(&response.LoadGuildName{
		CharacterID: characterID,
		GuildName:   "",
	}, broadcastOpt)
	ch.Broadcast(&response.LoadGuildIcon{
		CharacterID: characterID,
	}, broadcastOpt)
}

func (l *CharacterListenerImpl) OnMultiChat(ch *entity.Character, mode pconst.MultiChatMode, senderName string, message string) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.MultiChat{
		Mode:    mode,
		Name:    senderName,
		Message: message,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyStatusMessage(ch *entity.Character, code pconst.PartyStatusCode, name string) {
	if ch == nil {
		return
	}
	ch.Send(&response.PartyStatusMessage{
		Code: code,
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateJoin(ch *entity.Character, forChannel int32, partyID uint32, joinName string, leaderID uint32, members []response.PartyMemberStatus) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateJoin{
		ForChannel:           forChannel,
		PartyID:              partyID,
		JoiningCharacterName: joinName,
		LeaderCharacterID:    leaderID,
		Members:              members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateExpel(ch *entity.Character, forChannel int32, partyID uint32, targetID uint32, targetName string, leaderID uint32, members []response.PartyMemberStatus) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateExpel{
		ForChannel:          forChannel,
		PartyID:             partyID,
		TargetCharacterID:   targetID,
		TargetCharacterName: targetName,
		LeaderCharacterID:   leaderID,
		Members:             members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLeave(ch *entity.Character, forChannel int32, partyID uint32, targetID uint32, targetName string, leaderID uint32, members []response.PartyMemberStatus) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateLeave{
		ForChannel:          forChannel,
		PartyID:             partyID,
		TargetCharacterID:   targetID,
		TargetCharacterName: targetName,
		LeaderCharacterID:   leaderID,
		Members:             members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateDisband(ch *entity.Character, partyID uint32, leaderID uint32) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateDisband{
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLeaderChange(ch *entity.Character, newLeaderID uint32, byDisconnect bool) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateLeaderChange{
		NewLeaderCharacterID: newLeaderID,
		ByDisconnect:         byDisconnect,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLogOnOff(ch *entity.Character, forChannel int32, partyID uint32, leaderID uint32, members []response.PartyMemberStatus) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateLogOnOff{
		ForChannel:        forChannel,
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
		Members:           members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateSilent(ch *entity.Character, forChannel int32, partyID uint32, leaderID uint32, members []response.PartyMemberStatus) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.PartyUpdateSilent{
		ForChannel:        forChannel,
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
		Members:           members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyStatusMessage(ch *entity.Character, code pconst.BuddyStatusCode) {
	if ch == nil {
		return
	}
	ch.Send(&response.BuddyStatus{Code: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyListUpdate(ch *entity.Character, action pconst.BuddyListSyncAction, entries []response.BuddyEntry) {
	if ch == nil {
		return
	}
	ch.Send(&response.BuddyListUpdate{
		Action:  action,
		Entries: entries,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyChannelUpdate(ch *entity.Character, buddyCharacterID uint32, channel int32) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.BuddyChannelUpdate{
		CharacterID: buddyCharacterID,
		Channel:     channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyAddRequest(ch *entity.Character, fromCharacterID uint32, fromName string) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.BuddyAddRequest{
		FromCharacterID: fromCharacterID,
		FromName:        fromName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMessage(ch *entity.Character, code pconst.GuildResponseCode) {
	if ch == nil {
		return
	}
	_ = ch.Send(&response.GuildMessage{
		Code: code,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnExpGain(ch *entity.Character, exp uint32) {
	expPacket := &response.GainExp{
		Gain:  exp,
		White: false,
	}
	ch.Send(expPacket, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnControlMoveMob(ch *entity.Character, mob *entity.Mob, moveId uint16, enabledSkill bool, mp uint16, skillId uint32, skillLevel uint8) {
	ch.Send(&response.ControlMoveMob{
		OID:          mob.OID,
		MoveId:       moveId,
		EnabledSkill: enabledSkill,
		MP:           mp,
		SkillId:      uint8(skillId),
		SkillLevel:   skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMobHp(ch *entity.Character, mob *entity.Mob, percentage uint8) {
	ch.Send(&response.ShowMobHp{
		OID:        mob.OID,
		Percentage: percentage,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUnlockAction(ch *entity.Character) {
	ch.Send(&response.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnQuestStarted(ch *entity.Character, qp *entity.Quest, npcID uint32) {
	if ch == nil || qp == nil || qp.Wz == nil {
		return
	}
	ch.Send(&response.UpdateQuest{
		QuestStatus: dto.QuestStatus{
			QuestID:        uint16(qp.QuestID),
			Status:         uint8(qp.Status),
			MobKills:       qp.StartedMobKills(),
			Deadline:       qp.Deadline,
			StatusRecord:   qp.StatusRecord.AsString(),
			CompletionTime: qp.CompletionTime,
		},
	}, types.SEND_POLICY_ENCRYPT)
	if npcID != 0 {
		ch.Send(&response.UpdateQuestNPC{
			Progress:    8,
			QuestID:     uint16(qp.QuestID),
			NPCID:       npcID,
			NextQuestID: 0,
		}, types.SEND_POLICY_ENCRYPT)
	}
	if qp.Wz != nil && qp.Wz.Meta.TimeLimit > 0 {
		ch.Send(&response.Clock{
			Seconds: int32(qp.Wz.Meta.TimeLimit),
		}, types.SEND_POLICY_ENCRYPT)
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestCompleted(ch *entity.Character, qp *entity.Quest, npcID uint32, nextQuestID uint32) {
	if ch == nil || qp == nil || qp.Wz == nil {
		return
	}
	ch.Send(&response.UpdateQuest{
		QuestStatus: dto.QuestStatus{
			QuestID:        uint16(qp.QuestID),
			Status:         uint8(qp.Status),
			CompletionTime: qp.CompletionTime,
		},
	}, types.SEND_POLICY_ENCRYPT)
	ch.Send(&response.UpdateQuestNPC{
		Progress:    8,
		QuestID:     uint16(qp.QuestID),
		NPCID:       npcID,
		NextQuestID: nextQuestID,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestForfeited(ch *entity.Character, qp *entity.Quest) {
	if ch == nil || qp == nil || qp.Wz == nil {
		return
	}
	ch.Send(&response.UpdateQuest{
		QuestStatus: dto.QuestStatus{
			QuestID:        uint16(qp.QuestID),
			Status:         uint8(qp.Status),
			CompletionTime: qp.CompletionTime,
		},
	}, types.SEND_POLICY_ENCRYPT)
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestProgress(ch *entity.Character, qp *entity.Quest) {
	if ch == nil || qp == nil || qp.Wz == nil {
		return
	}
	ch.Send(&response.UpdateQuest{
		QuestStatus: dto.QuestStatus{
			QuestID:        uint16(qp.QuestID),
			Status:         response.QuestWireStatusStarted,
			MobKills:       qp.StartedMobKills(),
			Deadline:       qp.Deadline,
			StatusRecord:   qp.StatusRecord.AsString(),
			CompletionTime: qp.CompletionTime,
		},
	}, types.SEND_POLICY_ENCRYPT)

	if qp.CanComplete(ch, entity.QuestPhaseOpts{}) == nil {
		ch.Send(&response.ShowQuestCompletion{
			QuestID: uint16(qp.QuestID),
		}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnQuestRecordExChanged(ch *entity.Character, qp *entity.Quest) {
	if ch == nil || qp == nil {
		return
	}
	if qp.QuestID > 0xFFFF {
		return
	}
	ch.Send(&response.UpdateQuestRecordEx{
		QuestID: uint16(qp.QuestID),
		Data:    qp.RecordExWire(),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnItemGainFailed(ch *entity.Character, mode constant.ItemGainFailedType) {
	ch.Send(&response.ItemGainFailed{
		Mode: mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotUpdated(ch *entity.Character, inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)
	ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInventorySlotAdded(ch *entity.Character, inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)

	fromDrop := true
	if item != nil {
		model := item.GetModel()
		if model != nil && model.GetCapacity() >= 2 {
			fromDrop = false
		}
	}

	ch.Send(&response.AddInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
		FromDrop:      fromDrop,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowItemGain(ch *entity.Character, itemId uint32, count uint32, mode constant.ShowItemGainType) {
	ch.Send(&response.ShowItemGain{
		ItemId: itemId,
		Count:  count,
		Mode:   mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowMesoGain(ch *entity.Character, count int32, mode constant.ShowMesoGainType) {
	ch.Send(&response.ShowMesoGain{
		Count: count,
		Mode:  mode,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateStats(ch *entity.Character, stats map[constant.Stat]int32, unlock bool) {
	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: unlock,
	}, types.SEND_POLICY_ENCRYPT)
	if stats != nil {
		_, hasHP := stats[constant.StatHP]
		_, hasMaxHP := stats[constant.StatMaxHP]
		if hasHP || hasMaxHP {
			l.OnPartyMemberHPChanged(ch, nil)
		}
	}
}

func (l *CharacterListenerImpl) OnShowSelfSkillEffect(ch *entity.Character, effectType pconst.SkillEffectType, skillID uint32, skillLevel uint8, additional *uint8) {
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

func (l *CharacterListenerImpl) OnShowSkillEffect(ch *entity.Character, effectType pconst.SkillEffectType, skillID uint32, skillLevel uint8, additional *uint8) {
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

func (l *CharacterListenerImpl) OnShowSelfEffect(ch *entity.Character, effectType response.EffectType) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfEffect{
		Type: effectType,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowEffect(ch, effectType)
}

func (l *CharacterListenerImpl) OnShowEffect(ch *entity.Character, effectType response.EffectType) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowEffect{
		CharacterID: ch.GetID(),
		Type:        effectType,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfDragonBloodEffect(ch *entity.Character, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfDragonBloodEffect{
		SkillID:    skillID,
		SkillLevel: skillLevel,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowDragonBloodEffect(ch, skillID, skillLevel)
}

func (l *CharacterListenerImpl) OnShowDragonBloodEffect(ch *entity.Character, skillID uint32, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowDragonBloodEffect{
		CharacterID: ch.GetID(),
		SkillID:     skillID,
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfHPHealedEffect(ch *entity.Character, amount int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfHPHealedEffect{
		Amount: amount,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowHPHealedEffect(ch, amount)
}

func (l *CharacterListenerImpl) OnShowHPHealedEffect(ch *entity.Character, amount int32) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowHPHealedEffect{
		CharacterID: ch.GetID(),
		Amount:      amount,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfRewardItemAnimation(ch *entity.Character, itemID uint32, effect string) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfRewardItemAnimation{
		ItemID: itemID,
		Effect: effect,
	}, types.SEND_POLICY_ENCRYPT)
	l.OnShowRewardItemAnimation(ch, itemID, effect)
}

func (l *CharacterListenerImpl) OnShowRewardItemAnimation(ch *entity.Character, itemID uint32, effect string) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowRewardItemAnimation{
		CharacterID: ch.GetID(),
		ItemID:      itemID,
		Effect:      effect,
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfItemMakerSuccessEffect(ch *entity.Character) {
	if ch.GetMap() == nil {
		return
	}

	ch.Send(&response.ShowSelfItemMakerSuccessEffect{}, types.SEND_POLICY_ENCRYPT)
	l.OnShowItemMakerSuccessEffect(ch)
}

func (l *CharacterListenerImpl) OnShowItemMakerSuccessEffect(ch *entity.Character) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowItemMakerSuccessEffect{
		CharacterID: ch.GetID(),
	}, nil)
}

func (l *CharacterListenerImpl) OnShowSelfCraftingEffect(ch *entity.Character, effect string, time int32, mode int32) {
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

func (l *CharacterListenerImpl) OnShowCraftingEffect(ch *entity.Character, effect string, time int32, mode int32) {
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

func (l *CharacterListenerImpl) OnShowSelfDiceEffect(ch *entity.Character, effectID int32, skillID uint32, skillLevel uint8) {
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

func (l *CharacterListenerImpl) OnShowDiceEffect(ch *entity.Character, effectID int32, skillID uint32, skillLevel uint8) {
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

func (l *CharacterListenerImpl) OnShowScrollEffect(ch *entity.Character, success bool, destroyedByCurse bool) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(&response.ShowScrollEffect{
		CharacterID:      ch.GetID(),
		Success:          success,
		DestroyedByCurse: destroyedByCurse,
	}, &entity.ObjectBroadcastOption{
		WithMe: true,
	})
}

func (l *CharacterListenerImpl) OnScrolledItem(ch *entity.Character, scrollInventoryType constant.InventoryType, scrollSlot int16, scrollCount uint16, upgradedSlot int16, destroyed bool, potential bool, upgradedItem entity.Item) {
	upgradedItemDTO := entity.ItemToDTO(upgradedItem)
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

func (l *CharacterListenerImpl) OnMobMoved(ch *entity.Character, mob *entity.Mob, isAggroed bool, centerSplit int8, skill1 uint8, skill2 uint8, skill3 uint8, skill4 uint8, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	mapInstance := mob.GetMap()
	if mapInstance == nil {
		return
	}

	controllerTable := mapInstance.GetControllerTable()
	controller, exists := controllerTable.GetController(mob)

	movePacket := &response.MoveMob{
		IsAggroed:   isAggroed,
		CenterSplit: centerSplit,
		Skill1:      skill1,
		Skill2:      skill2,
		Skill3:      skill3,
		Skill4:      skill4,
		OID:         mob.OID,
		StartPoint:  startPoint,
		Movements:   movements,
	}

	if exists {
		controller.Broadcast(movePacket, nil)
	} else {
		mob.Broadcast(movePacket, nil)
	}
}

func (l *CharacterListenerImpl) OnPlayerMove(ch *entity.Character, startPoint types.Vector2[int16], fragments []dto.MoveFragment) {
	if ch.GetMap() == nil {
		return
	}

	characterDTO := ch.ToDTO()

	movePacket := &response.Move{
		Character:  characterDTO,
		Fragments:  fragments,
		StartPoint: startPoint,
	}

	ch.Broadcast(movePacket, nil)
}

func (l *CharacterListenerImpl) OnFieldRelocate(ch *entity.Character, spawnPoint uint8) {
	ch.Send(&response.FieldRelocate{Portal: spawnPoint}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) broadcastAttack(ch *entity.Character, packet types.Packet) {
	if ch.GetMap() == nil {
		return
	}

	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(ch, &response.Attack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	})
}

func (l *CharacterListenerImpl) OnRangedAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(ch, &response.RangedAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
		CashBullet:  0,
	})
}

func (l *CharacterListenerImpl) OnMagicAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	l.broadcastAttack(ch, &response.MagicAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	})
}

func (l *CharacterListenerImpl) OnEndSortInventory(ch *entity.Character, inventoryType constant.InventoryType) {
	ch.Send(&response.EndSortInventory{
		InventoryType: inventoryType,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSwapInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, source int16, dest int16, equipmentAction int8) {
	ch.Send(&response.SwapInventorySlot{
		InventoryType:   inventoryType,
		Source:          source,
		Dest:            dest,
		EquipmentAction: response.EquipmentActionType(equipmentAction),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnRemoveInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, slot int16) {
	ch.Send(&response.RemoveInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, slot int16, item entity.Item) {
	itemDTO := entity.ItemToDTO(item)
	ch.Send(&response.UpdateInventorySlot{
		InventoryType: inventoryType,
		Slot:          slot,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnFullMergeInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, source int16, dest int16, count uint16) {
	ch.Send(&response.FullMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		Count:         count,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartialMergeInventorySlot(ch *entity.Character, inventoryType constant.InventoryType, source int16, dest int16, sourceCount uint16, destCount uint16) {
	ch.Send(&response.PartialMergeInventorySlot{
		InventoryType: inventoryType,
		Source:        source,
		Dest:          dest,
		SourceCount:   sourceCount,
		DestCount:     destCount,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnUpdateCharacterLook(ch *entity.Character) {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	characterDTO := ch.ToDTO()

	lookPacket := &response.UpdateCharacterLook{
		Character: characterDTO,
	}

	ch.Broadcast(lookPacket, nil)
}

func (l *CharacterListenerImpl) OnNpcAction(ch *entity.Character, bytes []byte) {
	ch.Send(&response.NpcAction{
		Bytes: bytes,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnClassChange(ch *entity.Character, oldClass uint16, newClass uint16) {
	stats := map[constant.Stat]int32{
		constant.StatClass:       int32(newClass),
		constant.StatAvailableSP: int32(ch.SkillPoint),
	}

	ch.Send(&response.UpdateStats{
		Stats:        stats,
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
	l.gs.GetPartySystem().UpdateMemberAsync(nil, ch)
	l.gs.alliance.NotifyMemberFieldsChanged(ch)
}

func (l *CharacterListenerImpl) OnPartyMemberFieldsChanged(ch *entity.Character) {
	l.gs.GetPartySystem().UpdateMemberAsync(nil, ch)
	l.gs.alliance.NotifyMemberFieldsChanged(ch)
}

func (l *CharacterListenerImpl) OnPartyMemberHPChanged(ch *entity.Character, recipient *entity.Character) {
	if ch == nil {
		return
	}
	partyIDPtr := ch.GetPartyID()
	if partyIDPtr == nil {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	pkt := &response.UpdatePartyMemberHP{
		CharacterID: ch.GetID(),
		CurrentHP:   int32(ch.GetHp()),
		MaxHP:       int32(ch.GetMaxHp()),
	}
	if recipient != nil {
		if recipient.GetID() == ch.GetID() {
			return
		}
		if recipient.GetMap() != mapInstance {
			return
		}
		rPID := recipient.GetPartyID()
		if rPID == nil || *rPID != *partyIDPtr {
			return
		}
		_ = recipient.Send(pkt, types.SEND_POLICY_ENCRYPT)
		return
	}
	for _, obj := range mapInstance.GetObjects(constant.ObjectTypeCharacter) {
		peer, ok := obj.(*entity.Character)
		if !ok || peer == nil || peer.GetID() == ch.GetID() {
			continue
		}
		pid := peer.GetPartyID()
		if pid == nil || *pid != *partyIDPtr {
			continue
		}
		_ = peer.Send(pkt, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnBuffAdded(ch *entity.Character, buffID int32, remainingDuration time.Duration, values map[constant.BuffFlag]int32) {
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

func (l *CharacterListenerImpl) OnBuffRemoved(ch *entity.Character, flags []constant.BuffFlag) {
	ch.Send(&response.CancelSelfBuff{Buffs: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.CancelBuff{
		CharacterID: int32(ch.GetID()),
		Buffs:       flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnDebuffAdded(ch *entity.Character, debuff constant.DebuffFlag, x int16, skillID uint16, skillLevel uint16, durationMs int32) {
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

func (l *CharacterListenerImpl) OnDebuffRemoved(ch *entity.Character, flags []constant.DebuffFlag) {
	ch.Send(&response.RemoveSelfDebuff{Diseases: flags}, types.SEND_POLICY_ENCRYPT)

	ch.Broadcast(&response.RemoveDebuff{
		CharacterID: int32(ch.GetID()),
		Diseases:    flags,
	}, nil)
}

func (l *CharacterListenerImpl) OnSkillPassiveHook(ch *entity.Character, skillID uint32, hook string) {
	if ch == nil {
		return
	}
	CallPassiveSkillHook(ch, skillID, hook)
}

func (l *CharacterListenerImpl) OnUpdateSkill(ch *entity.Character, skillID uint32, level int32, masterLevel int32) {
	if ch == nil {
		return
	}
	ch.Send(&response.UpdateSkills{
		SkillID:     skillID,
		Level:       level,
		MasterLevel: masterLevel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSkillCooldown(ch *entity.Character, skillID uint32, remainingSec uint16) {
	ch.Send(&response.SkillCooldown{
		SkillID:      skillID,
		RemainingSec: uint32(remainingSec),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnHiddenChanged(ch *entity.Character, hidden bool) {
	ch.Send(&response.SuperHide{Hidden: hidden}, types.SEND_POLICY_ENCRYPT)

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}

	if hidden {
		ch.Broadcast(&response.LeavePlayer{ID: ch.GetID()}, &entity.ObjectBroadcastOption{
			RecipientsRoleBelowPivot: true,
		})
	} else {
		for _, obj := range mapInstance.GetObjects(constant.ObjectTypeCharacter) {
			if viewer, ok := obj.(*entity.Character); ok {
				ch.SendSpawnSyncToViewer(viewer)
			}
		}
	}
	if controllerTable := mapInstance.GetControllerTable(); controllerTable != nil {
		controllerTable.Update(ch)
	}
}

func (l *CharacterListenerImpl) OnSummonSpawn(ch *entity.Character, summon *entity.Summon) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SpawnSummon{
		OwnerID:      ch.GetID(),
		OID:          summon.OID,
		SkillID:      summon.SkillID,
		SkillLevel:   summon.SkillLevel,
		Position:     summon.Position,
		MovementType: summon.MovementType,
		SummonType:   summon.SummonType,
		Animated:     true,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonRemove(ch *entity.Character, summon *entity.Summon, animated bool) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.RemoveSummon{
		OwnerID:  ch.GetID(),
		OID:      summon.OID,
		Animated: animated,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonMove(ch *entity.Character, summon *entity.Summon, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.MoveSummon{
		CharacterID: ch.GetID(),
		OID:         summon.OID,
		StartPoint:  startPoint,
		Fragments:   movements,
	}
	ch.Broadcast(packet, nil)
}

func (l *CharacterListenerImpl) OnSummonAttack(ch *entity.Character, summon *entity.Summon, animation uint8, targets []entity.SummonAttackTarget) {
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
		SummonSkillID: uint32(summon.SkillID),
		Animation:     animation,
		Targets:       respTargets,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonSkill(ch *entity.Character, summon *entity.Summon, newStance uint8) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.SummonSkill{
		CharacterID: ch.GetID(),
		SummonOID:   summon.OID,
		NewStance:   newStance,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnSummonDamaged(ch *entity.Character, summon *entity.Summon, unknown uint8, damage uint32, monsterIdFrom uint32) {
	if summon.GetMap() == nil {
		return
	}
	packet := &response.DamageSummon{
		CharacterID:   ch.GetID(),
		SummonSkillID: uint32(summon.SkillID),
		Unknown:       unknown,
		Damage:        damage,
		MonsterIDFrom: monsterIdFrom,
	}
	ch.Broadcast(packet, &entity.ObjectBroadcastOption{WithMe: true})
}

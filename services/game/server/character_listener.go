package server

import (
	"context"
	"fmt"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
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

func (l *CharacterListenerImpl) OnDialogStyle(ch *entity.Character, npc uint32, message string, styles []uint32) {
	ch.Send(&response.DialogStyle{
		NPC:    npc,
		Text:   message,
		Styles: styles,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyCapacity(ch *entity.Character, capacity uint8) {
	ch.Send(&response.BuddyCapacityUpdate{Capacity: capacity}, types.SEND_POLICY_ENCRYPT)
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

func (l *CharacterListenerImpl) OnKeyMap(ch *entity.Character) {
	ch.Send(&response.KeyMap{Slots: ch.KeyLayout().Bindings()}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMessage(ch *entity.Character, messageType constant.ServerMessageType, message string) {
	l.OnNotice(ch, messageType, message, 0, false)
}

func (l *CharacterListenerImpl) OnNotice(ch *entity.Character, messageType constant.ServerMessageType, message string, channel int, ear bool) {
	ch.Send(&response.Notice{
		Message: message,
		Type:    messageType,
		Channel: channel,
		MegaEar: ear,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) BroadcastNoticeAsync(ctx actor.Context, ch *entity.Character, messageType constant.ServerMessageType, message string, ear bool) *async.Promise[*internal.BroadcastNoticeReply] {
	fail := func(err error) *async.Promise[*internal.BroadcastNoticeReply] {
		p := async.NewDeferred[*internal.BroadcastNoticeReply](ctx)
		p.SetError(err)
		return p
	}
	if l.gs == nil || l.gs.internalClient == nil || ctx == nil {
		return fail(fmt.Errorf("broadcast notice unavailable"))
	}
	channelID := l.gs.config.ChannelId
	worldID := l.gs.config.WorldId
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.BroadcastNoticeReply, error) {
			return l.gs.internalClient.BroadcastNotice(c, &internal.BroadcastNoticeRequest{
				WorldId:         worldID,
				SourceChannelId: channelID,
				MessageType:     uint32(messageType),
				Message:         message,
				MegaEar:         ear,
			})
		},
		func(reply *internal.BroadcastNoticeReply) error {
			if reply == nil || !reply.GetOk() {
				return fmt.Errorf("broadcast notice failed")
			}
			l.gs.BroadcastNotice(messageType, message, int(channelID)+1, ear)
			return nil
		},
	)
}

func (l *CharacterListenerImpl) OnClock(ch *entity.Character, seconds int32) {
	ch.Send(&response.Clock{
		Seconds: seconds,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStopClock(ch *entity.Character) {
	ch.Send(&response.StopClock{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyCreated(ch *entity.Character, partyID uint32) {
	ch.Send(&response.PartyCreated{
		PartyID: partyID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShowGuildInfo(ch *entity.Character) {
	if l.gs == nil {
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
	if l.gs == nil {
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
	if l.gs == nil {
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
	ch.Send(&response.PartyInvite{
		PartyID:     partyID,
		InviterName: inviterName,
		PartySearch: partySearch,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildInvite(ch *entity.Character, guildID uint32, inviterName string) {
	_ = ch.Send(&response.GuildInvite{
		GuildID:     guildID,
		InviterName: inviterName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceInvite(ch *entity.Character, inviterGuildID uint32, inviterName string, allianceName string) {
	_ = ch.Send(&response.AllianceInvite{
		InviterGuildID: inviterGuildID,
		InviterName:    inviterName,
		AllianceName:   allianceName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceCreate(ch *entity.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, membershipGuilds []dto.AllianceMembershipChangeGuild) {
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

func (l *CharacterListenerImpl) OnAllianceDisband(ch *entity.Character, allianceID uint32) {
	_ = ch.Send(&response.AllianceDisband{
		AllianceID: allianceID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceInfoBroadcast(ch *entity.Character, info *dto.AllianceInfo) {
	if info == nil {
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

func (l *CharacterListenerImpl) OnAllianceMemberRankChanged(ch *entity.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo) {
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

func (l *CharacterListenerImpl) OnAllianceGuildAdded(ch *entity.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, newGuildID uint32, addedGuild *dto.GuildInfo, members []dto.AllianceGuildMemberRank, joining bool, membershipGuild *dto.AllianceMembershipChangeGuild) {
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

func (l *CharacterListenerImpl) OnAllianceGuildLeft(ch *entity.Character, info *dto.AllianceInfo, removedGuildID uint32, removedGuild *dto.GuildInfo, removedMembers []dto.AllianceGuildMemberRank, expelled bool, leaving bool) {
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

func (l *CharacterListenerImpl) OnGuildNewMember(ch *entity.Character, guildID uint32, member dto.GuildMemberStatus) {
	_ = ch.Send(&response.GuildNewMember{
		GuildID: guildID,
		Member:  member,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildLeaveSelf(ch *entity.Character) {
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
	_ = ch.Send(&response.GuildMemberLeft{
		GuildID:     guildID,
		CharacterID: targetID,
		Name:        targetName,
		WasExpelled: wasExpelled,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildRankTitleChange(ch *entity.Character, guildID uint32, rankTitles [5]string) {
	_ = ch.Send(&response.GuildRankTitleChange{
		GuildID:    guildID,
		RankTitles: rankTitles,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberRankChange(ch *entity.Character, guildID uint32, targetID uint32, guildRank uint8) {
	_ = ch.Send(&response.GuildChangeRank{
		GuildID:     guildID,
		CharacterID: targetID,
		GuildRank:   guildRank,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildEmblemChange(ch *entity.Character, guildID uint32, logoBG uint16, logoBGColor uint8, logo uint16, logoColor uint8) {
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
	_ = ch.Send(&response.GuildNotice{
		GuildID: guildID,
		Notice:  notice,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildCapacityChange(ch *entity.Character, guildID uint32, capacity uint8) {
	_ = ch.Send(&response.GuildCapacityChange{
		GuildID:  guildID,
		Capacity: capacity,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildGPChange(ch *entity.Character, guildID uint32, gp uint32, level uint32, amount int32) {
	_ = ch.Send(&response.GuildUpdateGP{
		GuildID:    guildID,
		GP:         gp,
		GuildLevel: level,
	}, types.SEND_POLICY_ENCRYPT)
	if amount != 0 {
		_ = ch.Send(&response.ShowGPGain{Amount: amount}, types.SEND_POLICY_ENCRYPT)
	}
}

func (l *CharacterListenerImpl) OnGuildRanking(ch *entity.Character, npcID uint32, entries []dto.GuildRankingEntry) {
	_ = ch.Send(&response.GuildShowRanks{
		NPCID:   npcID,
		Entries: entries,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberOnlineChange(ch *entity.Character, guildID uint32, subjectCharacterID uint32, online bool) {
	_ = ch.Send(&response.GuildMemberOnline{
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Online:      online,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMemberFieldsChange(ch *entity.Character, guildID uint32, subjectCharacterID uint32, level uint32, classID uint32) {
	_ = ch.Send(&response.GuildMemberLevelClassUpdate{
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Level:       level,
		ClassID:     classID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberOnlineChange(ch *entity.Character, allianceID uint32, guildID uint32, subjectCharacterID uint32, online bool) {
	_ = ch.Send(&response.AllianceMemberOnline{
		AllianceID:  allianceID,
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Online:      online,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnAllianceMemberFieldsChange(ch *entity.Character, allianceID uint32, guildID uint32, subjectCharacterID uint32, level uint32, classID uint32) {
	_ = ch.Send(&response.AllianceUpdateMember{
		AllianceID:  allianceID,
		GuildID:     guildID,
		CharacterID: subjectCharacterID,
		Level:       level,
		ClassID:     classID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildDisbandSelf(ch *entity.Character, guildID uint32) {
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
	_ = ch.Send(&response.MultiChat{
		Mode:    mode,
		Name:    senderName,
		Message: message,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyStatusMessage(ch *entity.Character, code pconst.PartyStatusCode, name string) {
	ch.Send(&response.PartyStatusMessage{
		Code: code,
		Name: name,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateJoin(ch *entity.Character, forChannel int32, partyID uint32, joinName string, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateJoin{
		ForChannel:           forChannel,
		PartyID:              partyID,
		JoiningCharacterName: joinName,
		LeaderCharacterID:    leaderID,
		Members:              members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateExpel(ch *entity.Character, forChannel int32, partyID uint32, targetID uint32, targetName string, leaderID uint32, members []response.PartyMemberStatus) {
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
	_ = ch.Send(&response.PartyUpdateDisband{
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLeaderChange(ch *entity.Character, newLeaderID uint32, byDisconnect bool) {
	_ = ch.Send(&response.PartyUpdateLeaderChange{
		NewLeaderCharacterID: newLeaderID,
		ByDisconnect:         byDisconnect,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateLogOnOff(ch *entity.Character, forChannel int32, partyID uint32, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateLogOnOff{
		ForChannel:        forChannel,
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
		Members:           members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPartyUpdateSilent(ch *entity.Character, forChannel int32, partyID uint32, leaderID uint32, members []response.PartyMemberStatus) {
	_ = ch.Send(&response.PartyUpdateSilent{
		ForChannel:        forChannel,
		PartyID:           partyID,
		LeaderCharacterID: leaderID,
		Members:           members,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyStatusMessage(ch *entity.Character, code pconst.BuddyStatusCode) {
	ch.Send(&response.BuddyStatus{Code: code}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyListUpdate(ch *entity.Character, action pconst.BuddyListSyncAction, entries []response.BuddyEntry) {
	ch.Send(&response.BuddyListUpdate{
		Action:  action,
		Entries: entries,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyChannelUpdate(ch *entity.Character, buddyCharacterID uint32, channel int32) {
	_ = ch.Send(&response.BuddyChannelUpdate{
		CharacterID: buddyCharacterID,
		Channel:     channel,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnBuddyAddRequest(ch *entity.Character, fromCharacterID uint32, fromName string) {
	_ = ch.Send(&response.BuddyAddRequest{
		FromCharacterID: fromCharacterID,
		FromName:        fromName,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildMessage(ch *entity.Character, code pconst.GuildResponseCode) {
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

func (l *CharacterListenerImpl) OnOpenNpcShop(ch *entity.Character, shopID uint32, shop *wz.Shop) {
	ch.Send(&response.OpenNpcShop{
		ShopID: int32(shopID),
		Shop:   shop,
		Items:  l.gs.GetResources().Items,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) LoadStorageAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.LoadStorageReply] {
	accountID := ch.AccountID
	worldID := l.gs.config.WorldId
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.LoadStorageReply, error) {
			return l.gs.internalClient.LoadStorage(c, &internal.LoadStorageRequest{
				AccountId: accountID,
				WorldId:   worldID,
			})
		},
		nil,
	)
}

func (l *CharacterListenerImpl) storageTabs(storage *entity.Storage, invTypes ...constant.InventoryType) map[constant.InventoryType][]dto.Item {
	tabs := make(map[constant.InventoryType][]dto.Item, len(invTypes))
	for _, invType := range invTypes {
		items := make([]dto.Item, 0, len(storage.Tabs[invType]))
		for _, item := range storage.Tabs[invType] {
			items = append(items, item.ToDTO())
		}
		tabs[invType] = items
	}
	return tabs
}

func (l *CharacterListenerImpl) OnOpenStorage(ch *entity.Character, npcID uint32) {
	meso := ch.Storage.Meso
	ch.Send(&response.Storage{
		Result: pconst.StorageResultOpen,
		NpcID:  npcID,
		Slots:  ch.Storage.Slots,
		Meso:   &meso,
		Tabs: l.storageTabs(ch.Storage,
			constant.InventoryTypeEquipment,
			constant.InventoryTypeConsume,
			constant.InventoryTypeInstallation,
			constant.InventoryTypeETC,
			constant.InventoryTypeCash,
		),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageTabChanged(ch *entity.Character, result pconst.StorageResult, invType constant.InventoryType) {
	ch.Send(&response.Storage{
		Result: result,
		Slots:  ch.Storage.Slots,
		Tabs:   l.storageTabs(ch.Storage, invType),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageArranged(ch *entity.Character) {
	ch.Send(&response.Storage{
		Result: pconst.StorageResultArrange,
		Slots:  ch.Storage.Slots,
		Tabs: l.storageTabs(ch.Storage,
			constant.InventoryTypeEquipment,
			constant.InventoryTypeConsume,
			constant.InventoryTypeInstallation,
			constant.InventoryTypeETC,
			constant.InventoryTypeCash,
		),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageMesoChanged(ch *entity.Character) {
	meso := ch.Storage.Meso
	ch.Send(&response.Storage{
		Result: pconst.StorageResultMeso,
		Slots:  ch.Storage.Slots,
		Meso:   &meso,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStorageError(ch *entity.Character, result pconst.StorageResult) {
	ch.Send(&response.StorageError{
		Result: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) LoadParcelsAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.LoadParcelsReply] {
	req := &internal.LoadParcelsRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.LoadParcelsReply, error) {
			return l.gs.internalClient.LoadParcels(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) SendParcelAsync(ctx actor.Context, ch *entity.Character, recipient string, parcel *internal.Parcel, oneOfAKind bool, sender *internal.CharacterSaveEntry) *async.Promise[*internal.SendParcelReply] {
	req := &internal.SendParcelRequest{
		WorldId:         l.gs.config.WorldId,
		RecipientName:   recipient,
		Parcel:          parcel,
		SenderAccountId: sender.GetCharacter().GetAccountId(),
		Sender:          sender,
		OneOfAKind:      oneOfAKind,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.SendParcelReply, error) {
			return l.gs.internalClient.SendParcel(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) AddCashAsync(ctx actor.Context, ch *entity.Character, nxCash int32, maplePoint int32) *async.Promise[*internal.AddCashReply] {
	req := &internal.AddCashRequest{
		WorldId:    l.gs.config.WorldId,
		AccountId:  ch.AccountID,
		NxCash:     nxCash,
		MaplePoint: maplePoint,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.AddCashReply, error) {
			return l.gs.internalClient.AddCash(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) CreateCashCouponsAsync(ctx actor.Context, ch *entity.Character, kind internal.CashCouponKind, value uint32, count uint32) *async.Promise[*internal.CreateCashCouponsReply] {
	req := &internal.CreateCashCouponsRequest{
		WorldId: l.gs.config.WorldId,
		Kind:    kind,
		Value:   value,
		Count:   count,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CreateCashCouponsReply, error) {
			return l.gs.internalClient.CreateCashCoupons(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) ClaimParcelAsync(ctx actor.Context, ch *entity.Character, parcelID uint32) *async.Promise[*internal.ClaimParcelReply] {
	req := &internal.ClaimParcelRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ParcelId:    parcelID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.ClaimParcelReply, error) {
			return l.gs.internalClient.ClaimParcel(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) DeleteParcelAsync(ctx actor.Context, ch *entity.Character, parcelID uint32) *async.Promise[*internal.DeleteParcelReply] {
	req := &internal.DeleteParcelRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
		ParcelId:    parcelID,
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.DeleteParcelReply, error) {
			return l.gs.internalClient.DeleteParcel(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) CheckParcelArrivalsAsync(ctx actor.Context, ch *entity.Character) *async.Promise[*internal.CheckParcelArrivalsReply] {
	req := &internal.CheckParcelArrivalsRequest{
		WorldId:     l.gs.config.WorldId,
		CharacterId: ch.GetID(),
	}
	return async.NewTask(ctx, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CheckParcelArrivalsReply, error) {
			return l.gs.internalClient.CheckParcelArrivals(c, req)
		},
		nil,
	)
}

func (l *CharacterListenerImpl) parcelsDTO(parcels []*entity.Parcel) []dto.Parcel {
	out := make([]dto.Parcel, 0, len(parcels))
	for _, parcel := range parcels {
		p := dto.Parcel{
			ID:      parcel.ID,
			Sender:  parcel.Sender,
			Meso:    parcel.Meso,
			Expire:  parcel.ExpiresAt(),
			Quick:   parcel.Quick,
			Message: parcel.Message,
		}
		if parcel.Item != nil {
			p.Item = parcel.Item.ToDTO()
		}
		out = append(out, p)
	}
	return out
}

func (l *CharacterListenerImpl) OnOpenDuey(ch *entity.Character, fromArrival bool) {
	ch.Send(&response.DueyOpen{
		FromArrival: fromArrival,
		Parcels:     l.parcelsDTO(ch.Duey.Parcels),
		Expired:     l.parcelsDTO(ch.Duey.Expired),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyResult(ch *entity.Character, result pconst.DueyResult) {
	ch.Send(&response.Duey{
		Result: result,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyRemoved(ch *entity.Character, parcelID uint32, reason uint8) {
	ch.Send(&response.DueyRemoved{
		ParcelID: parcelID,
		Reason:   reason,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnDueyArrival(ch *entity.Character, sender string, quick bool, count int) {
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

func (l *CharacterListenerImpl) OnGuildBulletinThreadList(ch *entity.Character, threads []*internal.GuildBulletinBoardThreadEntry, start int, totalCount int, notice *internal.GuildBulletinBoardThreadEntry) {
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

func (l *CharacterListenerImpl) OnGuildBulletinThread(ch *entity.Character, detail *internal.GuildBulletinBoardThreadDetail) {
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

func (l *CharacterListenerImpl) OnShipState(ch *entity.Character, state uint16) {
	ch.Send(&response.ShipState{State: state}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnShipBalrog(ch *entity.Character) {
	ch.Send(&response.ShipSpecialEffect{Effect: response.ShipSpecialBalrog}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnSitOnChair(ch *entity.Character, itemID uint32) {
	ch.Broadcast(&response.ShowChair{
		CharacterID: ch.GetID(),
		ItemID:      itemID,
	}, nil)
}

func (l *CharacterListenerImpl) OnSitOnMapSeat(ch *entity.Character, seatID int16) {
	ch.Send(&response.CancelChair{
		ChairID: seatID,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnStandUp(ch *entity.Character) {
	ch.Send(&response.CancelChair{
		ChairID: -1,
	}, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(&response.ShowChair{
		CharacterID: ch.GetID(),
		ItemID:      0,
	}, nil)
}

func (l *CharacterListenerImpl) OnUnlockAction(ch *entity.Character) {
	ch.Send(&response.UpdateStats{
		UnlockAction: true,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnInspect(ch *entity.Character, profile *response.CharacterProfile) {
	ch.Send(profile, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCardRegistered(ch *entity.Character, cardID uint32, count uint32) {
	ch.Send(&response.MonsterBookSetCard{Success: true, CardID: cardID, Count: count}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCardFull(ch *entity.Character) {
	ch.Send(&response.MonsterBookSetCard{}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnMonsterBookCover(ch *entity.Character, cardID uint32) {
	ch.Send(&response.MonsterBookSetCover{CardID: cardID}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTeleportStones(ch *entity.Character, vip bool) {
	ch.Send(&response.TeleportStoneResult{Result: pconst.TeleportStoneResultList, VIP: vip, Maps: ch.TeleportStones.Slots(vip)}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnTeleportStoneFailed(ch *entity.Character, vip bool, result pconst.TeleportStoneResult) {
	ch.Send(&response.TeleportStoneResult{Result: result, VIP: vip}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnScriptError(ch *entity.Character, script string, err error) {
	if ch.GetRole() != constant.RoleAdmin {
		return
	}
	l.OnMessage(ch, constant.MsgPinkText, fmt.Sprintf("스크립트 오류 %s: %v", script, err))
}

func (l *CharacterListenerImpl) OnQuestStarted(ch *entity.Character, qp *entity.Quest, npcID uint32) {
	if qp == nil || qp.Wz == nil {
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
		l.OnClock(ch, int32(qp.Wz.Meta.TimeLimit))
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestCompleted(ch *entity.Character, qp *entity.Quest, npcID uint32, nextQuestID uint32) {
	if qp == nil || qp.Wz == nil {
		return
	}
	ch.Send(&response.UpdateQuest{
		QuestStatus: dto.QuestStatus{
			QuestID:        uint16(qp.QuestID),
			Status:         uint8(qp.Status),
			CompletionTime: qp.CompletionTime,
		},
	}, types.SEND_POLICY_ENCRYPT)
	if npcID != 0 {
		ch.Send(&response.UpdateQuestNPC{
			Progress:    8,
			QuestID:     uint16(qp.QuestID),
			NPCID:       npcID,
			NextQuestID: nextQuestID,
		}, types.SEND_POLICY_ENCRYPT)
	}
	l.OnUnlockAction(ch)
}

func (l *CharacterListenerImpl) OnQuestForfeited(ch *entity.Character, qp *entity.Quest) {
	if qp == nil {
		return
	}
	if qp.QuestID > 0xFFFF {
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

func (l *CharacterListenerImpl) OnShowQuestCompletion(ch *entity.Character, questID uint32) {
	if questID == 0 || questID > 0xFFFF {
		return
	}
	ch.Send(&response.ShowQuestCompletion{
		QuestID: uint16(questID),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPlaySound(ch *entity.Character, sound string, broadcast bool) {
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

func (l *CharacterListenerImpl) OnPlayPortalSound(ch *entity.Character) {
	ch.Send(&response.ShowSelfSkillEffect{
		Type:       pconst.SkillEffectTypePortal,
		SkillID:    0,
		SkillLevel: 1,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnQuestProgress(ch *entity.Character, qp *entity.Quest) {
	if qp == nil || qp.Wz == nil {
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
	if qp == nil {
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

	controllerID := uint32(0)
	if ch != nil {
		controllerID = ch.GetID()
	}
	for _, obj := range mob.Nears(constant.ObjectTypeCharacter, &entity.SearchOption{
		IncludeHidden: true,
	}) {
		peer, ok := obj.(*entity.Character)
		if !ok || peer == nil || peer.GetID() == controllerID {
			continue
		}
		_ = peer.Send(movePacket, types.SEND_POLICY_ENCRYPT)
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

func (l *CharacterListenerImpl) OnAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.Attack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	}, nil)
}

func (l *CharacterListenerImpl) OnRangedAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
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

func (l *CharacterListenerImpl) OnMagicAttack(ch *entity.Character, attackPayload dto.AttackPayload, skillLevel uint8) {
	if ch.GetMap() == nil {
		return
	}
	ch.Broadcast(&response.MagicAttack{
		AttackInfo:  attackPayload.ToAttackInfo(),
		CharacterId: ch.GetID(),
		SkillLevel:  skillLevel,
	}, nil)
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

	lookPacket := &response.UpdateCharacterLook{
		Character:      ch.ToDTO(),
		CrushRing:      ch.Inventory.WornRing(ch.Inventory.Rings.Left),
		FriendshipRing: ch.Inventory.WornRing(ch.Inventory.Rings.Mid),
		MarriageRing:   ch.Wedding.RingToDTO(),
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
		constant.StatAvailableSP: int32(ch.Points.SP),
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

func (l *CharacterListenerImpl) OnUpdateSkill(ch *entity.Character, skillID uint32, level int32, masterLevel int32) {
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
		for _, obj := range mapInstance.GetObjectsNear(ch.GetPosition(), constant.ObjectTypeCharacter, &entity.SearchOption{IncludeHidden: true}) {
			if viewer, ok := obj.(*entity.Character); ok {
				ch.SendSpawnSyncToViewer(viewer)
			}
		}
	}
	if controllerTable := mapInstance.GetControllerTable(); controllerTable != nil {
		controllerTable.Update(ch)
	}
}

func (l *CharacterListenerImpl) OnPetSpawn(ch *entity.Character) {
	ch.Broadcast(&response.SpawnPet{
		CharacterID: ch.GetID(),
		Pet:         ch.Pets.Active.ToDTO(),
	}, &entity.ObjectBroadcastOption{WithMe: true})
	l.OnPetExceptions(ch)
	sn := *ch.Pets.Active.Item.UniqueId
	ch.Send(&response.UpdateStats{Pet: &sn, UnlockAction: true}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetRemove(ch *entity.Character, reason constant.PetRemoveReason) {
	ch.Broadcast(&response.SpawnPet{
		CharacterID: ch.GetID(),
		Reason:      reason,
	}, &entity.ObjectBroadcastOption{WithMe: true})
	sn := uint64(0)
	ch.Send(&response.UpdateStats{Pet: &sn, UnlockAction: true}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetMove(ch *entity.Character, startPoint types.Vector2[int16], movements []dto.MoveFragment) {
	ch.Broadcast(&response.MovePet{
		CharacterID: ch.GetID(),
		StartPoint:  startPoint,
		Fragments:   movements,
	}, nil)
}

func (l *CharacterListenerImpl) OnPetChat(ch *entity.Character, typ uint8, action uint8, text string) {
	ch.Broadcast(&response.PetChat{
		CharacterID: ch.GetID(),
		Type:        typ,
		Action:      action,
		Text:        text,
	}, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetCommand(ch *entity.Character, index uint8, success bool) {
	ch.Broadcast(&response.PetCommand{CharacterID: ch.GetID(), Index: index, Success: success}, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetFood(ch *entity.Character, success bool) {
	ch.Broadcast(&response.PetCommand{CharacterID: ch.GetID(), Food: true, Success: success}, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetUpdated(ch *entity.Character, pet *entity.Pet) {
	slot, ok := ch.Inventory.FindSlot(constant.InventoryTypeCash, pet)
	if ok == false {
		return
	}
	l.OnInventorySlotUpdated(ch, constant.InventoryTypeCash, slot, pet)
}

func (l *CharacterListenerImpl) OnPetExceptions(ch *entity.Character) {
	ch.Send(&response.PetExceptions{
		CharacterID: ch.GetID(),
		SN:          *ch.Pets.Active.Item.UniqueId,
		ItemIDs:     ch.Pets.Active.Item.Exceptions,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetNameChanged(ch *entity.Character) {
	ch.Broadcast(&response.PetNameChanged{
		CharacterID: ch.GetID(),
		Name:        ch.Pets.Active.Item.Name,
	}, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnPetSkillChanged(ch *entity.Character, pet *entity.Pet, skill constant.PetSkill, add bool) {
	l.OnPetUpdated(ch, pet)
	ch.Send(&response.PetSkillChanged{SN: *pet.UniqueId, Add: add, Skill: skill}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnPetLevelUp(ch *entity.Character) {
	ch.Send(&response.ShowSelfEffect{Type: response.EffectTypePet}, types.SEND_POLICY_ENCRYPT)
	ch.Broadcast(&response.ShowEffect{CharacterID: ch.GetID(), Type: response.EffectTypePet}, nil)
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

func (l *CharacterListenerImpl) OnCarnivalStart(ch *entity.Character, team constant.CarnivalTeam, persAvail, persTotal, friendAvail, friendTotal, enemyAvail, enemyTotal int) {
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

func (l *CharacterListenerImpl) OnCarnivalObtainedCP(ch *entity.Character, avail, total int) {
	ch.Send(&response.CarnivalObtainedCP{
		Avail: uint16(avail),
		Total: uint16(total),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalPartyCP(ch *entity.Character, team constant.CarnivalTeam, avail, total int) {
	ch.Send(&response.CarnivalPartyCP{
		Team:  team,
		Avail: uint16(avail),
		Total: uint16(total),
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalSummon(ch *entity.Character, tab constant.CarnivalTab, num uint8, name string) {
	ch.Broadcast(&response.CarnivalSummon{
		Tab:  uint8(tab),
		Num:  num,
		Name: name,
	}, &entity.ObjectBroadcastOption{WithMe: true})
}

func (l *CharacterListenerImpl) OnCarnivalDied(ch *entity.Character, team constant.CarnivalTeam, name string, lostCP uint8) {
	ch.Send(&response.CarnivalDied{
		Team:   team,
		Name:   name,
		LostCP: lostCP,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnCarnivalResult(ch *entity.Character, winner bool) {
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

func (l *CharacterListenerImpl) OnConfirmShopTransaction(ch *entity.Character, code uint8) {
	_ = ch.Send(&response.ConfirmShopTransaction{
		Code: code,
	}, types.SEND_POLICY_ENCRYPT)
}

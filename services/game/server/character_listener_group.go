package server

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/types"
)

func (l *CharacterListenerImpl) OnBuddyCapacity(ch *game.Character, capacity uint8) {
	ch.Send(&response.BuddyCapacityUpdate{Capacity: capacity}, types.SEND_POLICY_ENCRYPT)
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

// OnBroadcastGuildAppearance shows the character's guild name and emblem over its head, or none once it has no guild.
func (l *CharacterListenerImpl) OnBroadcastGuildAppearance(ch *game.Character) {
	characterID := ch.GetID()
	namePkt := &response.LoadGuildName{CharacterID: characterID}
	iconPkt := &response.LoadGuildIcon{CharacterID: characterID}
	if guildID, ok := ch.Guild().ID(); ok {
		ent := l.gs.systems.Get[*game.GuildSystem]().Get(guildID)
		if ent == nil {
			return
		}
		namePkt.GuildName = ent.Name()
		if ent.Logo() != nil {
			iconPkt.HasGuildIcon = true
			iconPkt.LogoBG = uint16(ent.Logo().LogoBG())
			iconPkt.LogoBGColor = uint8(ent.Logo().LogoBGColor())
			iconPkt.Logo = uint16(ent.Logo().Logo())
			iconPkt.LogoColor = uint8(ent.Logo().LogoColor())
		}
	}
	broadcastOpt := &game.ObjectBroadcastOption{WithMe: true}
	ch.Broadcast(namePkt, broadcastOpt)
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

func (l *CharacterListenerImpl) OnAllianceGuildAdded(ch *game.Character, info *dto.AllianceInfo, guilds []*dto.GuildInfo, newGuildID uint32, addedGuild *dto.GuildInfo, members []dto.AllianceGuildMemberRank, joining bool, membershipGuilds []dto.AllianceMembershipChangeGuild) {
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
		if len(membershipGuilds) > 0 {
			_ = ch.Send(&response.AllianceChangeMembership{
				InAlliance: true,
				AllianceID: info.AllianceID,
				Guilds:     membershipGuilds,
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
	_ = ch.Send(&response.GuildShowInfo{
		Info: nil,
	}, types.SEND_POLICY_ENCRYPT)
}

func (l *CharacterListenerImpl) OnGuildExpelledSelf(ch *game.Character, guildID uint32) {
	_ = ch.Send(&response.GuildMemberLeft{
		GuildID:     guildID,
		CharacterID: ch.GetID(),
		Name:        ch.GetName(),
		WasExpelled: true,
	}, types.SEND_POLICY_ENCRYPT)
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
	_ = ch.Send(&response.GuildDisband{
		GuildID: guildID,
	}, types.SEND_POLICY_ENCRYPT)
	_ = ch.Send(&response.GuildShowInfo{
		Info: nil,
	}, types.SEND_POLICY_ENCRYPT)
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

func (l *CharacterListenerImpl) OnPartyMemberHP(recipient *game.Character, member *game.Character) {
	_ = recipient.Send(&response.UpdatePartyMemberHP{
		CharacterID: member.GetID(),
		CurrentHP:   int32(member.GetHp()),
		MaxHP:       int32(member.GetMaxHp()),
	}, types.SEND_POLICY_ENCRYPT)
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

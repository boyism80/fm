package game

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
)

type Alliance struct {
	gameWorld         GameWorld
	allianceID        uint32
	leaderCharacterID uint32
	name              string
	rankTitles        [5]string
	guildIDs          []uint32
	capacity          uint32
	notice            string
}

func (a *Alliance) Clone() *Alliance {
	if a == nil {
		return nil
	}
	return &Alliance{
		gameWorld:         a.gameWorld,
		allianceID:        a.allianceID,
		leaderCharacterID: a.leaderCharacterID,
		name:              a.name,
		rankTitles:        a.rankTitles,
		guildIDs:          append([]uint32(nil), a.guildIDs...),
		capacity:          a.capacity,
		notice:            a.notice,
	}
}

func (a *Alliance) BuildS2C() *dto.AllianceInfo {
	if a == nil {
		return nil
	}
	return &dto.AllianceInfo{
		AllianceID: a.allianceID,
		Name:       a.name,
		RankTitles: a.rankTitles,
		GuildIDs:   append([]uint32(nil), a.guildIDs...),
		Capacity:   a.capacity,
		Notice:     a.notice,
	}
}

func allianceFromProto(gw GameWorld, pb *internal.Alliance) *Alliance {
	if pb == nil || pb.GetAllianceId() == 0 {
		return nil
	}
	var rankTitles [5]string
	titles := pb.GetRankTitles()
	for i := 0; i < len(rankTitles) && i < len(titles); i++ {
		rankTitles[i] = titles[i]
	}
	return &Alliance{
		gameWorld:         gw,
		allianceID:        pb.GetAllianceId(),
		leaderCharacterID: pb.GetLeaderCharacterId(),
		name:              pb.GetName(),
		rankTitles:        rankTitles,
		guildIDs:          append([]uint32(nil), pb.GetGuildIds()...),
		capacity:          pb.GetCapacity(),
		notice:            pb.GetNotice(),
	}
}

func (a *Alliance) CanSendInvite(inviterCharacterID uint32) bool {
	if a == nil {
		return false
	}
	if uint32(len(a.guildIDs)) >= a.capacity {
		return false
	}
	return a.leaderCharacterID == inviterCharacterID
}

package entity

import "github.com/boyism80/fm/protocol/dto"

type Alliance struct {
	GameWorld         GameWorld
	AllianceID        uint32
	LeaderCharacterID uint32
	Name              string
	RankTitles        [5]string
	GuildIDs          []uint32
	Capacity          uint32
	Notice            string
}

func (a *Alliance) Clone() *Alliance {
	if a == nil {
		return nil
	}
	return &Alliance{
		GameWorld:         a.GameWorld,
		AllianceID:        a.AllianceID,
		LeaderCharacterID: a.LeaderCharacterID,
		Name:              a.Name,
		RankTitles:        a.RankTitles,
		GuildIDs:          append([]uint32(nil), a.GuildIDs...),
		Capacity:          a.Capacity,
		Notice:            a.Notice,
	}
}

func AllianceFromInfo(gw GameWorld, info *dto.AllianceInfo) *Alliance {
	if info == nil {
		return nil
	}
	guildIDs := append([]uint32(nil), info.GuildIDs...)
	return &Alliance{
		GameWorld:  gw,
		AllianceID: info.AllianceID,
		Name:       info.Name,
		RankTitles: info.RankTitles,
		GuildIDs:   guildIDs,
		Capacity:   info.Capacity,
		Notice:     info.Notice,
	}
}

func (a *Alliance) CanSendInvite(inviterCharacterID uint32) bool {
	if a == nil || inviterCharacterID == 0 {
		return false
	}
	if uint32(len(a.GuildIDs)) >= a.Capacity {
		return false
	}
	return a.LeaderCharacterID == inviterCharacterID
}

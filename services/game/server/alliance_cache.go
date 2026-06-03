package server

import (
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

type allianceCacheEntry struct {
	info     *dto.AllianceInfo
	guilds   []*dto.GuildInfo
	guildIDs []uint32
}

func (gs *GameServer) storeAllianceFromProto(alliancePb *internal.Alliance) {
	if gs == nil || alliancePb == nil || alliancePb.GetAllianceId() == 0 {
		return
	}
	info := entity.AllianceInfoFromProto(alliancePb)
	guilds := entity.AllianceCreateGuildsFromProto(alliancePb)
	guildIDs := append([]uint32(nil), alliancePb.GetGuildIds()...)
	gs.allianceMu.Lock()
	if gs.allianceCache == nil {
		gs.allianceCache = make(map[uint32]allianceCacheEntry)
	}
	gs.allianceCache[alliancePb.GetAllianceId()] = allianceCacheEntry{
		info:     info,
		guilds:   guilds,
		guildIDs: guildIDs,
	}
	gs.allianceMu.Unlock()
}

func (gs *GameServer) allianceCacheEntry(allianceID uint32) (allianceCacheEntry, bool) {
	if gs == nil || allianceID == 0 {
		return allianceCacheEntry{}, false
	}
	gs.allianceMu.RLock()
	entry, ok := gs.allianceCache[allianceID]
	gs.allianceMu.RUnlock()
	return entry, ok
}

func (gs *GameServer) allianceShowPackets(allianceID uint32) []types.Packet {
	entry, ok := gs.allianceCacheEntry(allianceID)
	if !ok || entry.info == nil {
		return []types.Packet{
			&response.AllianceShowInfo{Info: nil},
		}
	}
	return []types.Packet{
		&response.AllianceShowInfo{Info: entry.info},
		&response.AllianceShowGuilds{Guilds: entry.guilds},
	}
}

func (gs *GameServer) allianceShowPacketsForCharacter(ch *entity.Character) []types.Packet {
	if gs == nil || ch == nil {
		return nil
	}
	guildID, ok := ch.GetGuildID()
	if !ok || gs.guild == nil {
		return []types.Packet{
			&response.AllianceShowInfo{Info: nil},
		}
	}
	g := gs.guild.Get(guildID)
	if g == nil || g.AllianceID == 0 {
		return []types.Packet{
			&response.AllianceShowInfo{Info: nil},
		}
	}
	return gs.allianceShowPackets(g.AllianceID)
}

func (gs *GameServer) deliverAllianceMemberOnlineChange(allianceID uint32, subjectGuildID uint32, subjectID uint32, online bool) {
	if gs == nil || gs.guild == nil || allianceID == 0 || subjectID == 0 {
		return
	}
	guildIDs := gs.allianceGuildIDs(allianceID)
	if len(guildIDs) == 0 {
		if subjectGuildID == 0 {
			return
		}
		guildIDs = []uint32{subjectGuildID}
	}
	for _, gid := range guildIDs {
		g := gs.guild.Get(gid)
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 || memberID == subjectID {
				continue
			}
			if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
				continue
			}
			gs.EnsureSend(nil, memberID, &g_actor.DeliverAllianceMemberOnlineChange{
				CharacterID: memberID,
				AllianceID:  allianceID,
				GuildID:     subjectGuildID,
				SubjectID:   subjectID,
				Online:      online,
			})
		}
	}
}

func (gs *GameServer) allianceGuildIDs(allianceID uint32) []uint32 {
	entry, ok := gs.allianceCacheEntry(allianceID)
	if !ok || len(entry.guildIDs) == 0 {
		return nil
	}
	out := append([]uint32(nil), entry.guildIDs...)
	return out
}

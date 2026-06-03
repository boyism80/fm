package server

import (
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/types"
)

func (gs *GameServer) applyAllianceFromProto(alliancePb *internal.Alliance) {
	if gs == nil || gs.guild == nil || alliancePb == nil {
		return
	}
	for _, g := range alliancePb.GetGuilds() {
		ent := entity.GuildFromProto(g)
		if ent == nil {
			continue
		}
		gs.guild.Update(g)
	}
	gs.storeAllianceFromProto(alliancePb)
}

func allianceCreatePackets(alliancePb *internal.Alliance) []types.Packet {
	info := entity.AllianceInfoFromProto(alliancePb)
	guilds := entity.AllianceCreateGuildsFromProto(alliancePb)
	if info == nil {
		return nil
	}
	membership := entity.AllianceMembershipChangeGuildsFromProto(alliancePb)
	return []types.Packet{
		&response.AllianceCreate{
			Info:   info,
			Guilds: guilds,
		},
		&response.AllianceShowInfo{
			Info: info,
		},
		&response.AllianceShowGuilds{
			Guilds: guilds,
		},
		&response.AllianceChangeMembership{
			InAlliance: true,
			AllianceID: info.AllianceID,
			Guilds:     membership,
		},
	}
}

func (gs *GameServer) broadcastAllianceCreate(alliancePb *internal.Alliance) {
	if gs == nil || alliancePb == nil {
		return
	}
	packets := allianceCreatePackets(alliancePb)
	if len(packets) == 0 {
		return
	}
	for _, g := range alliancePb.GetGuilds() {
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 {
				continue
			}
			if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
				continue
			}
			gs.EnsureSend(nil, memberID, &g_actor.DeliverAllianceCreate{
				CharacterID: memberID,
				Packets:     packets,
			})
		}
	}
}

package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

const allianceSendOpcode uint16 = 0x31

func writeAllianceInfo(w *stream.StreamWriter, info *dto.AllianceInfo) {
	w.WriteU32(info.AllianceID)
	w.WriteStr16(info.Name)
	for i := 0; i < 5; i++ {
		w.WriteStr16(info.RankTitles[i])
	}
	w.WriteU8(uint8(len(info.GuildIDs)))
	for _, guildID := range info.GuildIDs {
		w.WriteU32(guildID)
	}
	w.WriteU32(info.Capacity)
	w.WriteStr16(info.Notice)
}

func writeAllianceMembershipChangeGuilds(w *stream.StreamWriter, inAlliance bool, guilds []dto.AllianceMembershipChangeGuild) {
	w.WriteU8(uint8(len(guilds)))
	for _, g := range guilds {
		w.WriteU32(g.GuildID)
		w.WriteU32(uint32(len(g.Members)))
		for _, m := range g.Members {
			w.WriteU32(m.CharacterID)
			if inAlliance {
				w.WriteU8(m.AllianceRank)
			} else {
				w.WriteU8(0)
			}
		}
	}
}

func writeAllianceGuildMemberRanks(w *stream.StreamWriter, added bool, members []dto.AllianceGuildMemberRank) {
	w.WriteU32(uint32(len(members)))
	for _, m := range members {
		w.WriteU32(m.CharacterID)
		if added {
			w.WriteU8(m.AllianceRank)
		} else {
			w.WriteU8(0)
		}
	}
}

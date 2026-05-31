package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

func writeGuildMemberData(w *stream.StreamWriter, members []dto.GuildMemberStatus) {
	w.WriteU8(uint8(len(members)))
	for _, member := range members {
		w.WriteU32(member.CharacterID)
	}
	for _, member := range members {
		w.WriteStaticStr(member.Name, 13)
		w.WriteU32(member.JobID)
		w.WriteU32(member.Level)
		w.WriteU32(member.GuildRank)
		online := uint32(0)
		if member.Online {
			online = 1
		}
		w.WriteU32(online)
		w.WriteU32(0)
		w.WriteU32(member.AllianceRank)
	}
}

func writeGuildInfo(w *stream.StreamWriter, info *dto.GuildInfo) {
	w.WriteU32(info.GuildID)
	w.WriteStr16(info.Name)
	for i := 0; i < 5; i++ {
		w.WriteStr16(info.RankTitles[i])
	}
	writeGuildMemberData(w, info.Members)
	w.WriteU32(info.Capacity)
	w.WriteU16(info.LogoBG)
	w.WriteU8(info.LogoBGColor)
	w.WriteU16(info.Logo)
	w.WriteU8(info.LogoColor)
	w.WriteStr16(info.Notice)
	w.WriteU32(info.GP)
	w.WriteU32(info.AllianceID)
}

func writeGuildMemberJoinedPayload(w *stream.StreamWriter, guildID uint32, member dto.GuildMemberStatus) {
	w.WriteU32(guildID)
	w.WriteU32(member.CharacterID)
	w.WriteStaticStr(member.Name, 13)
	w.WriteU32(member.JobID)
	w.WriteU32(member.Level)
	w.WriteU32(member.GuildRank)
	online := uint32(0)
	if member.Online {
		online = 1
	}
	w.WriteU32(online)
	w.WriteU32(1)
	w.WriteU32(member.AllianceRank)
}

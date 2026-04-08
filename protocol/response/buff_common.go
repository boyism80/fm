package response

import (
	"sort"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

func WriteBuff(writer *stream.StreamWriter, buff constant.BuffFlag) {
	pos0 := buff.Position - 1
	for i := 0; i < constant.MaxBuffFlag; i++ {
		var v int32
		if i == pos0 {
			v = int32(buff.Mask)
		}
		writer.Write32(v)
	}
}

func WriteBuffs(writer *stream.StreamWriter, buffs []constant.BuffFlag) {
	mask := make([]int32, constant.MaxBuffFlag)
	for _, buff := range buffs {
		idx := buff.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= int32(buff.Mask)
		}
	}
	for i := 0; i < constant.MaxBuffFlag; i++ {
		writer.Write32(mask[i])
	}
}

func WriteDebuff(writer *stream.StreamWriter, debuff constant.DebuffFlag) {
	pos0 := debuff.Position - 1
	for i := 0; i < constant.MaxBuffFlag; i++ {
		var v int32
		if i == pos0 {
			v = int32(debuff.Mask)
		}
		writer.Write32(v)
	}
}

func WriteDebuffs(writer *stream.StreamWriter, debuffs []constant.DebuffFlag) {
	mask := make([]int32, constant.MaxBuffFlag)
	for _, debuff := range debuffs {
		idx := debuff.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= int32(debuff.Mask)
		}
	}
	for i := 0; i < constant.MaxBuffFlag; i++ {
		writer.Write32(mask[i])
	}
}

func SortBuffEntries(buffs []dto.BuffEntry) {
	sort.Slice(buffs, func(i, j int) bool {
		if buffs[i].Buff.Position != buffs[j].Buff.Position {
			return buffs[i].Buff.Position < buffs[j].Buff.Position
		}
		return buffs[i].Buff.Mask < buffs[j].Buff.Mask
	})
}

package response

import (
	"sort"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

func writeBuffMask(writer *stream.StreamWriter, buffs []dto.BuffEntry) error {
	mask := make([]int32, constant.MaxBuffFlag)
	for _, e := range buffs {
		idx := e.Buff.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= int32(e.Buff.Mask)
		}
	}
	for i := 0; i < constant.MaxBuffFlag; i++ {
		if err := writer.Write32(mask[i]); err != nil {
			return err
		}
	}
	return nil
}

func writeMask(writer *stream.StreamWriter, buffs []constant.BuffFlag) error {
	mask := make([]int32, constant.MaxBuffFlag)
	for _, s := range buffs {
		idx := s.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= int32(s.Mask)
		}
	}
	for i := 0; i < constant.MaxBuffFlag; i++ {
		if err := writer.Write32(mask[i]); err != nil {
			return err
		}
	}
	return nil
}

func sortBuffs(buffs []dto.BuffEntry) {
	sort.Slice(buffs, func(i, j int) bool {
		if buffs[i].Buff.Position != buffs[j].Buff.Position {
			return buffs[i].Buff.Position < buffs[j].Buff.Position
		}
		return buffs[i].Buff.Mask < buffs[j].Buff.Mask
	})
}

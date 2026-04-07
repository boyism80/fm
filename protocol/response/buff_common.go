package response

import (
	"sort"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MaskSlot struct {
	Mask     uint32
	Position int
}

func WriteSingleMask(writer *stream.StreamWriter, slot MaskSlot) error {
	pos0 := slot.Position - 1
	for i := 0; i < constant.MaxBuffFlag; i++ {
		var v int32
		if i == pos0 {
			v = int32(slot.Mask)
		}
		writer.Write32(v)
	}
	return nil
}

func WriteMask(writer *stream.StreamWriter, slots []MaskSlot) error {
	mask := make([]int32, constant.MaxBuffFlag)
	for _, s := range slots {
		idx := s.Position - 1
		if idx >= 0 && idx < constant.MaxBuffFlag {
			mask[idx] |= int32(s.Mask)
		}
	}
	for i := 0; i < constant.MaxBuffFlag; i++ {
		writer.Write32(mask[i])
	}
	return nil
}

func SortBuffEntries(buffs []dto.BuffEntry) {
	sort.Slice(buffs, func(i, j int) bool {
		if buffs[i].Buff.Position != buffs[j].Buff.Position {
			return buffs[i].Buff.Position < buffs[j].Buff.Position
		}
		return buffs[i].Buff.Mask < buffs[j].Buff.Mask
	})
}

func SlotsFromBuffEntries(buffs []dto.BuffEntry) []MaskSlot {
	slots := make([]MaskSlot, len(buffs))
	for i := range buffs {
		slots[i] = MaskSlot{Mask: buffs[i].Buff.Mask, Position: buffs[i].Buff.Position}
	}
	return slots
}

func SlotsFromBuffFlags(buffs []constant.BuffFlag) []MaskSlot {
	slots := make([]MaskSlot, len(buffs))
	for i := range buffs {
		slots[i] = MaskSlot{Mask: buffs[i].Mask, Position: buffs[i].Position}
	}
	return slots
}

func SlotFromDebuffFlag(d constant.DebuffFlag) MaskSlot {
	return MaskSlot{Mask: d.Mask, Position: d.Position}
}

func SlotsFromDebuffFlags(debuffs []constant.DebuffFlag) []MaskSlot {
	slots := make([]MaskSlot, len(debuffs))
	for i := range debuffs {
		slots[i] = SlotFromDebuffFlag(debuffs[i])
	}
	return slots
}

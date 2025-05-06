package entity

import (
	"math"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
)

func (item *Item) Serialize(writer *stream.StreamWriter, zeroPosition, leaveOut, trade, bagSlot bool) {
	parts := item.Parts
	if zeroPosition {
		if !leaveOut {
			writer.WriteU8(0)
		}
	} else {
		if parts <= -1 {
			parts *= -1
			if parts > 100 && parts < 1000 {
				parts -= 100
			}
		}
		if bagSlot {
			writer.WriteU32(uint32((parts % 100) - 1))
		} else if !trade && item.Type == 1 { // equipment
			writer.WriteU16(uint16(parts))
		} else {
			writer.WriteU8(uint8(parts))
		}
	}

	if item.Pet != nil {
		writer.WriteU8(3)
	} else {
		writer.WriteU8(item.Type)
	}
	writer.WriteU32(item.Id)

	hasUID := item.UniqueID > 0
	writer.WriteBoolean(hasUID)
	if hasUID {
		writer.Write64(item.UniqueID)
	}

	if item.Pet != nil {
		_ = item.serializePet(writer)
	} else {
		writer.WriteU64(util.GetTime(item.Expiration))

		isEquipment := (item.Type == 1 && item.Equip != nil)
		if isEquipment {
			equip := item.Equip
			writer.WriteU8(equip.UpgradeSlots)
			writer.WriteU8(equip.Level)
			writer.WriteU16(equip.Str)
			writer.WriteU16(equip.Dex)
			writer.WriteU16(equip.Int)
			writer.WriteU16(equip.Luk)
			writer.WriteU16(equip.Hp)
			writer.WriteU16(equip.Mp)
			writer.WriteU16(equip.Watk)
			writer.WriteU16(equip.Matk)
			writer.WriteU16(equip.Wdef)
			writer.WriteU16(equip.Mdef)
			writer.WriteU16(equip.Acc)
			writer.WriteU16(equip.Avoid)
			writer.WriteU16(equip.Hands)
			writer.WriteU16(equip.Speed)
			writer.WriteU16(equip.Jump)
			writer.WriteStr16(item.Owner)
			writer.WriteU16(item.Flag)
			writer.WriteBoolean(item.IncSkill > 0)

			writer.WriteU8(uint8(math.Max(float64(equip.BaseLevel), float64(equip.EquipLevel))))
			writer.WriteU32(equip.ExpPercentage * 100000)

			if item.UniqueID <= 0 {
				if item.InventoryID <= 0 {
					writer.Write64(-1)
				} else {
					writer.WriteU64(item.InventoryID)
				}
			}

			writer.WriteU64(util.GetTime(-2))
			writer.Write32(-1)
		} else {
			writer.WriteU16(item.Quantity)
			writer.WriteStr16(item.Owner)
			writer.WriteU16(item.Flag)

			if item.IsThrowingStar() || item.IsBullet() || item.Id/10000 == 287 {
				if item.InventoryID <= 0 {
					writer.WriteU64(0xFFFFFFFFFFFFFFFF)
				} else {
					writer.WriteU64(item.InventoryID)
				}
			}
		}
	}
}

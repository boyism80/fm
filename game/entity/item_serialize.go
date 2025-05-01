package entity

import (
	"errors"
	"math"
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/common/util"
)

func (item *Item) serializePet(writer *stream.StreamWriter) error {
	if item.Pet == nil {
		return errors.New("serializePet: pet is nil")
	}

	writer.WriteU64(util.GetTime(item.Expiration))
	writer.WriteStaticStr(item.Name, 13)
	writer.WriteU8(item.Pet.Level)
	writer.WriteU16(item.Pet.Closeness)
	writer.WriteU8(item.Pet.Fullness)

	if item == nil {
		writer.WriteU64(util.GetKoreanTimestamp(int64(float64(time.Now().UnixNano()) / float64(time.Millisecond) * 1.5)))
	} else {
		writer.WriteU64(util.GetTime(item.Expiration))
	}

	writer.WriteU16(item.Pet.Speed)
	writer.WriteU16(item.Pet.Flags)

	if item.Pet.PetItemId == 5000054 && item.Pet.SecondsLeft > 0 {
		writer.WriteU32(item.Pet.SecondsLeft)
	} else {
		writer.WriteU32(0)
	}

	return nil
}

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
		} else if !trade && item.Type == 1 {
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

	hasUniqueId := item.UniqueID > 0
	writer.WriteU8(func() uint8 {
		if hasUniqueId {
			return 1
		}
		return 0
	}())

	if hasUniqueId {
		writer.WriteU64(item.UniqueID)
	}

	if item.Pet != nil {
		_ = item.serializePet(writer)
	} else {
		writer.WriteU64(util.GetTime(item.Expiration))

		if item.Type == 1 && item.Equip != nil {
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
			writer.WriteU8(func() uint8 {
				if equip.IncSkill > 0 {
					return 1
				}
				return 0
			}())

			writer.WriteU8(uint8(math.Max(float64(equip.BaseLevel), float64(equip.EquipLevel))))
			writer.WriteU32(equip.ExpPercentage * 100000)

			if item.UniqueID <= 0 {
				if item.InventoryID <= 0 {
					writer.WriteU64(0xFFFFFFFFFFFFFFFF)
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

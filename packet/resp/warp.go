package resp

import (
	"github.com/boyism80/fm/entity"
	"github.com/boyism80/fm/stream"
)

type Warp struct {
	Character *entity.Character
}

func (a *Warp) Serialize(writer *stream.StreamWriter) error {
	a.Character.Inventory[entity.InventoryTypeCash].SlotLimit = 60

	writer.WriteU16(0x55)
	writer.WriteU32(0) // channel
	writer.WriteU8(0)
	writer.WriteU8(1) // first time

	isEvent := false
	if isEvent {
		writer.WriteU16(1)
		writer.WriteStr16("event alarm")
		writer.WriteStr16("event message")
	} else {
		writer.WriteU16(0)
	}

	a.Character.Random1.Serialize(writer)

	writer.WriteU64(0xFFFFFFFFFFFFFFFF) // flag

	// flag 0x1
	a.Character.SerializeStats(writer)
	writer.WriteU8(20) // buddy capacity

	// flag 0x2 ~ 0x40
	a.Character.SerializeInventory(writer)

	// flag 0x100
	a.Character.SerializeSkills(writer)

	// flag 0x8000
	a.Character.SerializeCooldowns(writer)

	// flag 0x200, 0x4000
	a.Character.SerializeQuests(writer)

	// flag 0x400, 0x800
	a.Character.SerializeRings(writer)

	// flag 0x1000
	a.Character.SerializeRocks(writer)

	// flag 0x20000, 0x10000
	a.Character.SerializeMonsterBook(writer)

	// flag 0x40000
	a.Character.SerializeQuestInfo(writer)

	writer.WriteU16(0)
	writer.WriteU64(133906314709330000)
	return nil
}

func (a *Warp) Deserialize(reader *stream.StreamReader) error {
	return nil
}

package resp

import (
	"time"

	"github.com/boyism80/fm/dao"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/util"
)

type Warp struct {
}

func (a *Warp) Serialize(writer *stream.StreamWriter) error {
	ch := dao.NewDummyCharacter(1, "채승현")

	writer.WriteU16(0x55)
	writer.WriteU32(0) // channel
	writer.WriteU8(0)
	writer.WriteU8(0) // first time

	isEvent := true
	if isEvent {
		writer.WriteU16(1)
		writer.WriteStr16("event alarm")
		writer.WriteStr16("event message")
	} else {
		writer.WriteU16(0)
	}

	ch.Random1.Serialize(writer)

	writer.WriteU64(0xFFFFFFFFFFFFFFFF) // flag

	// flag 0x1
	ch.SerializeStats(writer)
	writer.WriteU8(20) // buddy capacity

	// flag 0x2 ~ 0x40
	ch.SerializeInventory(writer)

	// flag 0x100
	ch.SerializeSkills(writer)

	// flag 0x8000
	ch.SerializeCooldowns(writer)

	// flag 0x200, 0x4000
	ch.SerializeQuests(writer)

	// flag 0x400, 0x800
	ch.SerializeRings(writer)

	// flag 0x1000
	ch.SerializeRocks(writer)

	// flag 0x20000, 0x10000
	ch.SerializeMonsterBook(writer)

	// flag 0x40000
	ch.SerializeQuestInfo(writer)

	writer.WriteU16(0)
	writer.WriteU64(util.GetTime(time.Now().UnixMilli()))
	return nil
}

func (a *Warp) Deserialize(reader *stream.StreamReader) error {
	return nil
}

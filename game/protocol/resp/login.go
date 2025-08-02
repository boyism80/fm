package resp

import (
	"time"

	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
)

type Login struct {
	Character *entity.Character
}

func (p *Login) Serialize(writer *stream.StreamWriter) error {
	p.Character.Inventory[constant.INVENTORY_TYPE_CASH].SlotLimit = 60

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

	p.Character.Random1.Serialize(writer)
	p.Character.Serialize(writer)
	writer.WriteDateTime(time.Now())
	return nil
}

func (a *Login) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (l *Login) Opcode() uint16 {
	return 0x55
}

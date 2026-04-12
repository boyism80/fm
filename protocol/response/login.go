package response

import (
	"time"

	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type Login struct {
	Character *dto.Character
}

func (p *Login) Serialize(writer *stream.StreamWriter) error {

	if p.Character.Inventory != nil {
		if cashInv := p.Character.Inventory[constant.INVENTORY_TYPE_CASH]; cashInv != nil {
			cashInv.SlotLimit = 60
		}
	}

	writer.WriteU32(0)
	writer.WriteU8(0)
	writer.WriteU8(1)

	isEvent := false
	if isEvent {
		writer.WriteU16(1)
		writer.WriteStr16("event alarm")
		writer.WriteStr16("event message")
	} else {
		writer.WriteU16(0)
	}

	if p.Character.Random1 != nil {
		p.Character.Random1.Serialize(writer)
	}

	p.Character.Serialize(writer)
	writer.WriteDateTime(time.Now())
	return nil
}

func (a *Login) Deserialize(reader *stream.StreamReader) {
}

func (l *Login) Opcode() uint16 {
	return 0x55
}

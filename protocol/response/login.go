package response

import (
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type Login struct {
	Channel uint32
	CharacterInfo
}

func (p *Login) Serialize(writer *stream.StreamWriter) error {

	if p.Character.Inventory != nil {
		if cashInv := p.Character.Inventory.Tabs[constant.InventoryTypeCash]; cashInv != nil {
			cashInv.SlotLimit = 60
		}
	}

	writer.WriteU32(p.Channel)
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

	p.CharacterInfo.Serialize(writer)
	writer.WriteDateTime(clock.Now())
	return nil
}

func (a *Login) Deserialize(reader *stream.StreamReader) {
	a.Channel = reader.ReadU32()
	reader.ReadU8()
	reader.ReadU8()
	if reader.ReadU16() != 0 {
		return
	}
	reader.Skip(12)
	a.CharacterInfo.Deserialize(reader)
}

func (l *Login) Opcode() uint16 {
	return 0x55
}

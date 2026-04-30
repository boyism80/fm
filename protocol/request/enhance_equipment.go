package request

import (
	"github.com/boyism80/fm/stream"
)

type EnhanceEquipment struct {
	Tick       uint32
	ScrollSlot int16
	TargetSlot int16
}

func (*EnhanceEquipment) Opcode() byte { return 0x45 }

func (p *EnhanceEquipment) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *EnhanceEquipment) Deserialize(reader *stream.StreamReader) {
	p.Tick = reader.ReadU32()
	p.ScrollSlot = reader.Read16()
	p.TargetSlot = reader.Read16()
}

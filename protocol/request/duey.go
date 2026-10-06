package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type Duey struct {
	Mode          constant.DueyMode
	Value         int32
	Kind          int32
	InventoryType uint8
	Slot          int16
	Count         int16
	Meso          int32
	Recipient     string
	Quick         bool
	Message       string
	Coupon        int32
	ParcelID      uint32
}

func (*Duey) Opcode() byte { return 0x30 }

func (p *Duey) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case constant.DueyOpenFromArrival, constant.DueyIdentity:
		writer.Write32(p.Value)
		writer.Write32(p.Kind)
	case constant.DueySend:
		writer.WriteU8(p.InventoryType)
		writer.Write16(p.Slot)
		writer.Write16(p.Count)
		writer.Write32(p.Meso)
		writer.WriteStr16(p.Recipient)
		writer.WriteBoolean(p.Quick)
		if p.Quick {
			writer.WriteStr16(p.Message)
			writer.Write32(p.Coupon)
		}
	case constant.DueyReceive, constant.DueyDelete:
		writer.WriteU32(p.ParcelID)
	}
	return nil
}

func (p *Duey) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.DueyMode(reader.ReadU8())
	switch p.Mode {
	case constant.DueyOpenFromArrival, constant.DueyIdentity:
		p.Value = reader.Read32()
		p.Kind = reader.Read32()
	case constant.DueySend:
		p.InventoryType = reader.ReadU8()
		p.Slot = reader.Read16()
		p.Count = reader.Read16()
		p.Meso = reader.Read32()
		p.Recipient = reader.ReadStr16()
		p.Quick = reader.ReadBool()
		if p.Quick {
			p.Message = reader.ReadStr16()
			p.Coupon = reader.Read32()
		}
	case constant.DueyReceive, constant.DueyDelete:
		p.ParcelID = reader.ReadU32()
	}
}

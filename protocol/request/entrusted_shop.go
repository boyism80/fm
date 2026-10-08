package request

import (
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/stream"
)

type UseEntrustedShop struct{}

func (*UseEntrustedShop) Opcode() byte { return 0x2E }

func (p *UseEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *UseEntrustedShop) Deserialize(reader *stream.StreamReader) {}

type MiniRoom struct {
	Mode          constant.MiniRoomMode
	Type          uint8
	Title         string
	Slot          int16
	ItemID        uint32
	SN            uint32
	Message       string
	InventoryType uint8
	Bundles       uint16
	PerBundle     uint16
	Price         int32
	Index         uint16
}

func (*MiniRoom) Opcode() byte { return 0x65 }

func (p *MiniRoom) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case constant.MiniRoomCreate:
		writer.WriteU8(p.Type)
		writer.WriteStr16(p.Title)
		writer.WriteU8(0)
		writer.Write16(p.Slot)
		writer.WriteU32(p.ItemID)
	case constant.MiniRoomVisit:
		writer.WriteU32(p.SN)
		writer.WriteU8(0)
	case constant.MiniRoomChat:
		writer.WriteStr16(p.Message)
	case constant.MiniRoomOpen:
		writer.WriteU8(1)
	case constant.MiniRoomAddItem:
		writer.WriteU8(p.InventoryType)
		writer.Write16(p.Slot)
		writer.WriteU16(p.Bundles)
		writer.WriteU16(p.PerBundle)
		writer.Write32(p.Price)
	case constant.MiniRoomBuy:
		writer.WriteU8(uint8(p.Index))
		writer.WriteU16(p.Bundles)
	case constant.MiniRoomRemoveItem:
		writer.WriteU16(p.Index)
	}
	return nil
}

func (p *MiniRoom) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.MiniRoomMode(reader.ReadU8())
	switch p.Mode {
	case constant.MiniRoomCreate:
		p.Type = reader.ReadU8()
		p.Title = reader.ReadStr16()
		reader.ReadU8()
		p.Slot = reader.Read16()
		p.ItemID = reader.ReadU32()
	case constant.MiniRoomVisit:
		p.SN = reader.ReadU32()
	case constant.MiniRoomChat:
		p.Message = reader.ReadStr16()
	case constant.MiniRoomAddItem:
		p.InventoryType = reader.ReadU8()
		p.Slot = reader.Read16()
		p.Bundles = reader.ReadU16()
		p.PerBundle = reader.ReadU16()
		p.Price = reader.Read32()
	case constant.MiniRoomBuy:
		p.Index = uint16(reader.ReadU8())
		p.Bundles = reader.ReadU16()
	case constant.MiniRoomRemoveItem:
		p.Index = reader.ReadU16()
	}
}

type StoreBank struct {
	Mode constant.StoreBankMode
}

func (*StoreBank) Opcode() byte { return 0x2F }

func (p *StoreBank) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	return nil
}

func (p *StoreBank) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.StoreBankMode(reader.ReadU8())
}

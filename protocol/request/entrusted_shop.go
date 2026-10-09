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
	Name          string
	Names         []string
	TargetID      uint32
	Reason        uint8
	Count         uint16
	TradeSlot     uint8
	Meso          int32
}

func (*MiniRoom) Opcode() byte { return 0x65 }

func (p *MiniRoom) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case constant.MiniRoomCreate:
		writer.WriteU8(p.Type)
		if p.Type == constant.MiniRoomTypeTrade {
			writer.WriteU8(0)
			break
		}
		writer.WriteStr16(p.Title)
		writer.WriteU8(0)
		writer.Write16(p.Slot)
		writer.WriteU32(p.ItemID)
	case constant.MiniRoomInvite:
		writer.WriteU32(p.TargetID)
	case constant.MiniRoomDecline:
		writer.WriteU32(p.SN)
		writer.WriteU8(p.Reason)
	case constant.MiniRoomVisit:
		writer.WriteU32(p.SN)
		writer.WriteU8(0)
	case constant.MiniRoomTradePutItem:
		writer.WriteU8(p.InventoryType)
		writer.Write16(p.Slot)
		writer.WriteU16(p.Count)
		writer.WriteU8(p.TradeSlot)
	case constant.MiniRoomTradePutMeso:
		writer.Write32(p.Meso)
	case constant.MiniRoomChat:
		writer.WriteStr16(p.Message)
	case constant.MiniRoomOpen:
		writer.WriteU8(1)
	case constant.MiniRoomAddItem, constant.MiniRoomPersonalAddItem:
		writer.WriteU8(p.InventoryType)
		writer.Write16(p.Slot)
		writer.WriteU16(p.Bundles)
		writer.WriteU16(p.PerBundle)
		writer.Write32(p.Price)
	case constant.MiniRoomBuy, constant.MiniRoomPersonalBuy:
		writer.WriteU8(uint8(p.Index))
		writer.WriteU16(p.Bundles)
	case constant.MiniRoomRemoveItem, constant.MiniRoomPersonalRemoveItem:
		writer.WriteU16(p.Index)
	case constant.MiniRoomKick, constant.MiniRoomKickTimeout:
		writer.WriteU8(uint8(p.Slot))
		writer.WriteStr16(p.Name)
	case constant.MiniRoomBlacklist:
		writer.WriteU16(uint16(len(p.Names)))
		for _, name := range p.Names {
			writer.WriteStr16(name)
		}
	}
	return nil
}

func (p *MiniRoom) Deserialize(reader *stream.StreamReader) {
	p.Mode = constant.MiniRoomMode(reader.ReadU8())
	switch p.Mode {
	case constant.MiniRoomCreate:
		p.Type = reader.ReadU8()
		if p.Type == constant.MiniRoomTypeTrade {
			reader.ReadU8()
			break
		}
		p.Title = reader.ReadStr16()
		reader.ReadU8()
		p.Slot = reader.Read16()
		p.ItemID = reader.ReadU32()
	case constant.MiniRoomInvite:
		p.TargetID = reader.ReadU32()
	case constant.MiniRoomDecline:
		p.SN = reader.ReadU32()
		p.Reason = reader.ReadU8()
	case constant.MiniRoomVisit:
		p.SN = reader.ReadU32()
	case constant.MiniRoomTradePutItem:
		p.InventoryType = reader.ReadU8()
		p.Slot = reader.Read16()
		p.Count = reader.ReadU16()
		p.TradeSlot = reader.ReadU8()
	case constant.MiniRoomTradePutMeso:
		p.Meso = reader.Read32()
	case constant.MiniRoomChat:
		p.Message = reader.ReadStr16()
	case constant.MiniRoomAddItem, constant.MiniRoomPersonalAddItem:
		p.InventoryType = reader.ReadU8()
		p.Slot = reader.Read16()
		p.Bundles = reader.ReadU16()
		p.PerBundle = reader.ReadU16()
		p.Price = reader.Read32()
	case constant.MiniRoomBuy, constant.MiniRoomPersonalBuy:
		p.Index = uint16(reader.ReadU8())
		p.Bundles = reader.ReadU16()
	case constant.MiniRoomRemoveItem, constant.MiniRoomPersonalRemoveItem:
		p.Index = reader.ReadU16()
	case constant.MiniRoomKick, constant.MiniRoomKickTimeout:
		p.Slot = int16(reader.ReadU8())
		p.Name = reader.ReadStr16()
	case constant.MiniRoomBlacklist:
		count := int(reader.ReadU16())
		p.Names = make([]string, 0, count)
		for i := 0; i < count; i++ {
			p.Names = append(p.Names, reader.ReadStr16())
		}
	}
}

type RemoteEntrustedShop struct {
	Slot int16
}

func (*RemoteEntrustedShop) Opcode() byte { return 0x2A }

func (p *RemoteEntrustedShop) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(p.Slot)
	return nil
}

func (p *RemoteEntrustedShop) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.Read16()
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

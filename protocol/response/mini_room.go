package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MiniRoomItem struct {
	Bundles   uint16
	PerBundle uint16
	Price     int32
	Item      dto.Item
}

type MiniRoomSale struct {
	ItemID  uint32
	Bundles uint16
	Total   int32
	Buyer   string
}

type MiniRoomVisitor struct {
	Slot      uint8
	Character *dto.Character
}

type MiniRoomItems struct {
	Meso  int32
	Items []MiniRoomItem
}

func (p *MiniRoomItems) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomItems) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultItems))
	p.serializeItems(writer)
	return nil
}

func (p *MiniRoomItems) serializeItems(writer *stream.StreamWriter) {
	writer.Write32(p.Meso)
	writer.WriteU8(uint8(len(p.Items)))
	for _, item := range p.Items {
		writer.WriteU16(item.Bundles)
		writer.WriteU16(item.PerBundle)
		writer.Write32(item.Price)
		item.Item.Serialize(writer, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
	}
}

func (p *MiniRoomItems) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.deserializeItems(reader)
}

func (p *MiniRoomItems) deserializeItems(reader *stream.StreamReader) {
	p.Meso = reader.Read32()
	count := int(reader.ReadU8())
	p.Items = make([]MiniRoomItem, 0, count)
	for i := 0; i < count; i++ {
		p.Items = append(p.Items, MiniRoomItem{
			Bundles:   reader.ReadU16(),
			PerBundle: reader.ReadU16(),
			Price:     reader.Read32(),
			Item:      dto.NewItemFromStream(reader),
		})
	}
}

type MiniRoomEnter struct {
	MySlot       uint8
	PermitItemID uint32
	OwnerName    string
	Visitors     []MiniRoomVisitor
	Elapsed      uint32
	FirstTime    bool
	Sold         []MiniRoomSale
	Title        string
	MaxItems     uint8
	MiniRoomItems
}

func (p *MiniRoomEnter) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomEnter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultEnter))
	writer.WriteU8(pconst.MiniRoomTypeHiredMerchant)
	writer.WriteU8(pconst.MiniRoomHiredMerchantUsers)
	writer.WriteU8(p.MySlot)
	writer.WriteU8(0)
	writer.WriteU32(p.PermitItemID)
	writer.WriteStr16(p.OwnerName)
	for _, visitor := range p.Visitors {
		writer.WriteU8(visitor.Slot)
		visitor.Character.SerializeLook(writer)
		writer.WriteStr16(visitor.Character.Name)
	}
	writer.WriteU8(0xFF)
	writer.WriteU16(0)
	writer.WriteStr16(p.OwnerName)
	if p.MySlot == 0 {
		writer.WriteU32(p.Elapsed)
		writer.WriteBoolean(p.FirstTime)
		writer.WriteU8(uint8(len(p.Sold)))
		for _, sale := range p.Sold {
			writer.WriteU32(sale.ItemID)
			writer.WriteU16(sale.Bundles)
			writer.Write32(sale.Total)
			writer.WriteStr16(sale.Buyer)
		}
		writer.Write32(p.Meso)
	}
	writer.WriteStr16(p.Title)
	writer.WriteU8(p.MaxItems)
	p.serializeItems(writer)
	return nil
}

func (p *MiniRoomEnter) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	reader.ReadU8()
	p.MySlot = reader.ReadU8()
	for slot := int8(reader.ReadU8()); slot >= 0; slot = int8(reader.ReadU8()) {
		if slot == 0 {
			p.PermitItemID = reader.ReadU32()
			p.OwnerName = reader.ReadStr16()
			continue
		}
		visitor := MiniRoomVisitor{Slot: uint8(slot), Character: &dto.Character{}}
		visitor.Character.DeserializeLook(reader)
		visitor.Character.Name = reader.ReadStr16()
		p.Visitors = append(p.Visitors, visitor)
	}
	for i := reader.ReadU16(); i > 0; i-- {
		reader.ReadStr16()
		reader.ReadU8()
	}
	p.OwnerName = reader.ReadStr16()
	if p.MySlot == 0 {
		p.Elapsed = reader.ReadU32()
		p.FirstTime = reader.ReadBool()
		for i := reader.ReadU8(); i > 0; i-- {
			p.Sold = append(p.Sold, MiniRoomSale{
				ItemID:  reader.ReadU32(),
				Bundles: reader.ReadU16(),
				Total:   reader.Read32(),
				Buyer:   reader.ReadStr16(),
			})
		}
		reader.Read32()
	}
	p.Title = reader.ReadStr16()
	p.MaxItems = reader.ReadU8()
	p.deserializeItems(reader)
}

type MiniRoomEnterFailed struct {
	Error pconst.MiniRoomEnterError
}

func (p *MiniRoomEnterFailed) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomEnterFailed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultEnter))
	writer.WriteU8(0)
	writer.WriteU8(uint8(p.Error))
	return nil
}

func (p *MiniRoomEnterFailed) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	p.Error = pconst.MiniRoomEnterError(reader.ReadU8())
}

type MiniRoomVisited struct {
	MiniRoomVisitor
}

func (p *MiniRoomVisited) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomVisited) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultVisitor))
	writer.WriteU8(p.Slot)
	p.Character.SerializeLook(writer)
	writer.WriteStr16(p.Character.Name)
	return nil
}

func (p *MiniRoomVisited) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Slot = reader.ReadU8()
	p.Character = &dto.Character{}
	p.Character.DeserializeLook(reader)
	p.Character.Name = reader.ReadStr16()
}

type MiniRoomChat struct {
	Slot    uint8
	Message string
}

func (p *MiniRoomChat) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultChat))
	writer.WriteU8(pconst.MiniRoomChatShop)
	writer.WriteU8(p.Slot)
	writer.WriteStr16(p.Message)
	return nil
}

func (p *MiniRoomChat) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	p.Slot = reader.ReadU8()
	p.Message = reader.ReadStr16()
}

type MiniRoomLeave struct {
	Slot   uint8
	Reason pconst.MiniRoomLeaveReason
}

func (p *MiniRoomLeave) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomLeave) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultLeave))
	writer.WriteU8(p.Slot)
	writer.WriteU8(uint8(p.Reason))
	return nil
}

func (p *MiniRoomLeave) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Slot = reader.ReadU8()
	p.Reason = pconst.MiniRoomLeaveReason(reader.ReadU8())
}

type MiniRoomBuyFailed struct {
	Result pconst.MiniRoomBuyResult
}

func (p *MiniRoomBuyFailed) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomBuyFailed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultBuy))
	writer.WriteU8(uint8(p.Result))
	return nil
}

func (p *MiniRoomBuyFailed) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Result = pconst.MiniRoomBuyResult(reader.ReadU8())
}

type MiniRoomArranged struct {
	Meso int32
}

func (p *MiniRoomArranged) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomArranged) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultArranged))
	writer.Write32(p.Meso)
	return nil
}

func (p *MiniRoomArranged) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Meso = reader.Read32()
}

type MiniRoomClosed struct {
	Result pconst.MiniRoomCloseResult
}

func (p *MiniRoomClosed) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomClosed) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultClosed))
	writer.WriteU8(uint8(p.Result))
	return nil
}

func (p *MiniRoomClosed) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Result = pconst.MiniRoomCloseResult(reader.ReadU8())
}

type MiniRoomMesoWithdrawn struct{}

func (p *MiniRoomMesoWithdrawn) Opcode() uint16 {
	return 0xEF
}

func (p *MiniRoomMesoWithdrawn) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultMesoWithdrawn))
	return nil
}

func (p *MiniRoomMesoWithdrawn) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
}

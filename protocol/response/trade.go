package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type TradeInvite struct {
	Inviter string
	SN      uint32
}

func (p *TradeInvite) Opcode() uint16 {
	return 0xEF
}

func (p *TradeInvite) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultInvite))
	writer.WriteU8(pconst.MiniRoomTypeTrade)
	writer.WriteStr16(p.Inviter)
	writer.WriteU32(p.SN)
	return nil
}

func (p *TradeInvite) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	p.Inviter = reader.ReadStr16()
	p.SN = reader.ReadU32()
}

type TradeInviteResult struct {
	Result pconst.MiniRoomInviteResult
	Name   string
}

func (p *TradeInviteResult) Opcode() uint16 {
	return 0xEF
}

func (p *TradeInviteResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultInviteResult))
	writer.WriteU8(uint8(p.Result))
	if p.Result != pconst.MiniRoomInviteNotFound {
		writer.WriteStr16(p.Name)
	}
	return nil
}

func (p *TradeInviteResult) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Result = pconst.MiniRoomInviteResult(reader.ReadU8())
	if p.Result != pconst.MiniRoomInviteNotFound {
		p.Name = reader.ReadStr16()
	}
}

type TradeEnter struct {
	MySlot  uint8
	Members []MiniRoomVisitor
}

func (p *TradeEnter) Opcode() uint16 {
	return 0xEF
}

func (p *TradeEnter) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultEnter))
	writer.WriteU8(pconst.MiniRoomTypeTrade)
	writer.WriteU8(pconst.MiniRoomTradeUsers)
	writer.WriteU8(p.MySlot)
	for _, member := range p.Members {
		writer.WriteU8(member.Slot)
		member.Character.SerializeLook(writer)
		writer.WriteStr16(member.Character.Name)
	}
	writer.WriteU8(0xFF)
	return nil
}

func (p *TradeEnter) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU8()
	reader.ReadU8()
	p.MySlot = reader.ReadU8()
	for slot := int8(reader.ReadU8()); slot >= 0; slot = int8(reader.ReadU8()) {
		member := MiniRoomVisitor{Slot: uint8(slot), Character: &dto.Character{}}
		member.Character.DeserializeLook(reader)
		member.Character.Name = reader.ReadStr16()
		p.Members = append(p.Members, member)
	}
}

type TradeItem struct {
	Who  uint8
	Slot uint8
	Item dto.Item
}

func (p *TradeItem) Opcode() uint16 {
	return 0xEF
}

func (p *TradeItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultTradeItem))
	writer.WriteU8(p.Who)
	writer.WriteU8(p.Slot)
	p.Item.Serialize(writer, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
	return nil
}

func (p *TradeItem) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Who = reader.ReadU8()
	p.Slot = reader.ReadU8()
	p.Item = dto.NewItemFromStream(reader)
}

type TradeMeso struct {
	Who  uint8
	Meso int32
}

func (p *TradeMeso) Opcode() uint16 {
	return 0xEF
}

func (p *TradeMeso) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultTradeMeso))
	writer.WriteU8(p.Who)
	writer.Write32(p.Meso)
	return nil
}

func (p *TradeMeso) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Who = reader.ReadU8()
	p.Meso = reader.Read32()
}

type TradeConfirm struct{}

func (p *TradeConfirm) Opcode() uint16 {
	return 0xEF
}

func (p *TradeConfirm) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.MiniRoomResultTradeConfirm))
	return nil
}

func (p *TradeConfirm) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
}

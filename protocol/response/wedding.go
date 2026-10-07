package response

import (
	"github.com/boyism80/fm/stream"
)

type Weather struct {
	ItemID  uint32
	Message string
}

func (p *Weather) Opcode() uint16 {
	return 0x60
}

func (p *Weather) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.ItemID)
	if p.ItemID != 0 {
		writer.WriteStr16(p.Message)
	}
	return nil
}

func (p *Weather) Deserialize(reader *stream.StreamReader) {
	p.ItemID = reader.ReadU32()
	if p.ItemID != 0 {
		p.Message = reader.ReadStr16()
	}
}

type YellowChat struct {
	Message string
}

func (p *YellowChat) Opcode() uint16 {
	return 0x3C
}

func (p *YellowChat) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(0xFF)
	writer.WriteStr16(p.Message)
	return nil
}

func (p *YellowChat) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Message = reader.ReadStr16()
}

type WeddingCouple struct {
	GroomID uint32
	BrideID uint32
}

func (p *WeddingCouple) Opcode() uint16 {
	return 0xF5
}

func (p *WeddingCouple) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.GroomID)
	writer.WriteU32(p.BrideID)
	return nil
}

func (p *WeddingCouple) Deserialize(reader *stream.StreamReader) {
	p.GroomID = reader.ReadU32()
	p.BrideID = reader.ReadU32()
}

type WeddingEffect struct {
}

func (p *WeddingEffect) Opcode() uint16 {
	return 0xF6
}

func (p *WeddingEffect) Serialize(writer *stream.StreamWriter) error {
	return nil
}

func (p *WeddingEffect) Deserialize(reader *stream.StreamReader) {
}

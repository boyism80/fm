package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type EngageRequest struct {
	Mode        pconst.EngageRequestMode
	Name        string
	CharacterID uint32
}

func (p *EngageRequest) Opcode() uint16 {
	return 0x37
}

func (p *EngageRequest) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	if p.Mode == pconst.EngageRequestPropose {
		writer.WriteStr16(p.Name)
		writer.WriteU32(p.CharacterID)
	}
	return nil
}

func (p *EngageRequest) Deserialize(reader *stream.StreamReader) {
	p.Mode = pconst.EngageRequestMode(reader.ReadU8())
	if p.Mode == pconst.EngageRequestPropose {
		p.Name = reader.ReadStr16()
		p.CharacterID = reader.ReadU32()
	}
}

type EngageResult struct {
	Result      pconst.EngageResult
	Marriage    *dto.Marriage
	GroomName   string
	BrideName   string
	WeddingType uint16
}

func (p *EngageResult) Opcode() uint16 {
	return 0x38
}

func (p *EngageResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	switch p.Result {
	case pconst.EngageResultEngaged, pconst.EngageResultMarried:
		p.Marriage.Serialize(writer)
	case pconst.EngageResultInvitation:
		writer.WriteStr16(p.GroomName)
		writer.WriteStr16(p.BrideName)
		writer.WriteU16(p.WeddingType)
	}
	return nil
}

func (p *EngageResult) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.EngageResult(reader.ReadU8())
	switch p.Result {
	case pconst.EngageResultEngaged, pconst.EngageResultMarried:
		p.Marriage = &dto.Marriage{
			ID:          reader.ReadU32(),
			GroomID:     reader.ReadU32(),
			BrideID:     reader.ReadU32(),
			Status:      reader.ReadU16(),
			GroomItemID: reader.ReadU32(),
			BrideItemID: reader.ReadU32(),
			GroomName:   reader.ReadStaticStr(13),
			BrideName:   reader.ReadStaticStr(13),
		}
	case pconst.EngageResultInvitation:
		p.GroomName = reader.ReadStr16()
		p.BrideName = reader.ReadStr16()
		p.WeddingType = reader.ReadU16()
	}
}

type SpouseMap struct {
	MapID       uint32
	CharacterID uint32
}

func (p *SpouseMap) Opcode() uint16 {
	return 0x3A
}

func (p *SpouseMap) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.MapID)
	writer.WriteU32(p.CharacterID)
	return nil
}

func (p *SpouseMap) Deserialize(reader *stream.StreamReader) {
	p.MapID = reader.ReadU32()
	p.CharacterID = reader.ReadU32()
}

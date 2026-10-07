package dto

import "github.com/boyism80/fm/stream"

type Marriage struct {
	ID          uint32
	GroomID     uint32
	BrideID     uint32
	Status      uint16
	GroomItemID uint32
	BrideItemID uint32
	GroomName   string
	BrideName   string
}

func (m *Marriage) Serialize(writer *stream.StreamWriter) {
	writer.WriteU32(m.ID)
	writer.WriteU32(m.GroomID)
	writer.WriteU32(m.BrideID)
	writer.WriteU16(m.Status)
	writer.WriteU32(m.GroomItemID)
	writer.WriteU32(m.BrideItemID)
	writer.WriteStaticStr(m.GroomName, 13)
	writer.WriteStaticStr(m.BrideName, 13)
}

func (m *Marriage) Deserialize(reader *stream.StreamReader) {
	m.ID = reader.ReadU32()
	m.GroomID = reader.ReadU32()
	m.BrideID = reader.ReadU32()
	m.Status = reader.ReadU16()
	m.GroomItemID = reader.ReadU32()
	m.BrideItemID = reader.ReadU32()
	m.GroomName = reader.ReadStaticStr(13)
	m.BrideName = reader.ReadStaticStr(13)
}

type MarriageRing struct {
	CharacterID uint32
	PartnerID   uint32
	ItemID      uint32
}

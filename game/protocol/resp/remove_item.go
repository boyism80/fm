package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type RemoveItemType uint8

const (
	RemoveItemTypeExpired RemoveItemType = iota
	RemoveItemTypeNoAnimated
	RemoveItemTypeAnimated
	RemoveItemTypeExplosion
	RemoveItemTypeLootByPet
)

type RemoveItem struct {
	Mode        RemoveItemType
	OID         uint32
	CharacterId uint32
}

func (p *RemoveItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0xC7)
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU32(p.OID)
	switch p.Mode {
	case RemoveItemTypeAnimated, RemoveItemTypeLootByPet:
		writer.WriteU32(p.CharacterId)
	}
	return nil
}

func (p *RemoveItem) Deserialize(reader *stream.StreamReader) error {
	return nil
}

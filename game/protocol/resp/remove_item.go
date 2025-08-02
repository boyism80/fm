package resp

import (
	"github.com/boyism80/fm/common/stream"
)

type RemoveItemType uint8

const (
	REMOVE_ITEM_TYPE_EXPIRED RemoveItemType = iota
	REMOVE_ITEM_TYPE_NO_ANIMATED
	REMOVE_ITEM_TYPE_ANIMATED
	REMOVE_ITEM_TYPE_EXPLOSION
	REMOVE_ITEM_TYPE_LOOT_BY_PET
)

type RemoveItem struct {
	Mode        RemoveItemType
	OID         uint32
	CharacterId uint32
}

func (p *RemoveItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU32(p.OID)
	switch p.Mode {
	case REMOVE_ITEM_TYPE_ANIMATED, REMOVE_ITEM_TYPE_LOOT_BY_PET:
		writer.WriteU32(p.CharacterId)
	}
	return nil
}

func (p *RemoveItem) Deserialize(reader *stream.StreamReader) error {
	return nil
}

func (p *RemoveItem) Opcode() uint16 {
	return 0xC7
}

package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type RemoveItem struct {
	Mode        constant.RemoveItemType
	OID         uint32
	CharacterId uint32
}

func (p *RemoveItem) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	writer.WriteU32(p.OID)
	switch p.Mode {
	case constant.REMOVE_ITEM_TYPE_EXPLOSION:
		writer.Write16(655)
	case constant.REMOVE_ITEM_TYPE_ANIMATED, constant.REMOVE_ITEM_TYPE_LOOT_BY_PET:
		writer.WriteU32(p.CharacterId)
	}
	return nil
}

func (p *RemoveItem) Deserialize(reader *stream.StreamReader) {
}

func (p *RemoveItem) Opcode() uint16 {
	return 0xC7
}

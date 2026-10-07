package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type WeddingGiftTabs map[constant.InventoryType][]dto.Item

func (tabs WeddingGiftTabs) serialize(writer *stream.StreamWriter) {
	var mask uint64
	for typ := range tabs {
		mask |= 4 << (uint(typ) - 1)
	}
	writer.WriteU64(mask)
	for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
		items, ok := tabs[typ]
		if ok == false {
			continue
		}
		writer.WriteU8(uint8(len(items)))
		for _, item := range items {
			item.Serialize(writer, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
		}
	}
}

func readWeddingGiftTabs(reader *stream.StreamReader) WeddingGiftTabs {
	mask := reader.ReadU64()
	tabs := WeddingGiftTabs{}
	for typ := constant.InventoryTypeEquipment; typ <= constant.InventoryTypeCash; typ++ {
		if mask&(4<<(uint(typ)-1)) == 0 {
			continue
		}
		count := reader.ReadU8()
		items := make([]dto.Item, 0, count)
		for range count {
			items = append(items, dto.NewItemFromStream(reader))
		}
		tabs[typ] = items
	}
	return tabs
}

type WeddingGift struct {
	Mode   pconst.WeddingGiftMode
	Wishes []string
	Tabs   WeddingGiftTabs
}

func (p *WeddingGift) Opcode() uint16 {
	return 0x39
}

func (p *WeddingGift) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Mode))
	switch p.Mode {
	case pconst.WeddingGiftOpenGive:
		writer.WriteU8(uint8(len(p.Wishes)))
		for _, wish := range p.Wishes {
			writer.WriteStr16(wish)
		}
	case pconst.WeddingGiftGiven:
		writer.WriteU8(uint8(len(p.Wishes)))
		for _, wish := range p.Wishes {
			writer.WriteStr16(wish)
		}
		p.Tabs.serialize(writer)
	case pconst.WeddingGiftOpenReceive, pconst.WeddingGiftReceived:
		p.Tabs.serialize(writer)
	}
	return nil
}

func (p *WeddingGift) Deserialize(reader *stream.StreamReader) {
	p.Mode = pconst.WeddingGiftMode(reader.ReadU8())
	switch p.Mode {
	case pconst.WeddingGiftOpenGive, pconst.WeddingGiftGiven:
		count := reader.ReadU8()
		p.Wishes = make([]string, 0, count)
		for range count {
			p.Wishes = append(p.Wishes, reader.ReadStr16())
		}
		if p.Mode == pconst.WeddingGiftGiven {
			p.Tabs = readWeddingGiftTabs(reader)
		}
	case pconst.WeddingGiftOpenReceive, pconst.WeddingGiftReceived:
		p.Tabs = readWeddingGiftTabs(reader)
	}
}

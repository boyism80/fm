package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type Storage struct {
	Result pconst.StorageResult
	NpcID  uint32
	Slots  uint8
	Meso   *int32
	Tabs   map[constant.InventoryType][]dto.Item
}

func (p *Storage) Opcode() uint16 {
	return 0xEA
}

func (p *Storage) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	if p.Result == pconst.StorageResultOpen {
		writer.WriteU32(p.NpcID)
	}
	writer.WriteU8(p.Slots)

	var mask uint64
	if p.Meso != nil {
		mask |= 2
	}
	for invType := range p.Tabs {
		mask |= 2 << invType
	}
	writer.WriteU64(mask)

	if p.Meso != nil {
		writer.Write32(*p.Meso)
	}
	for invType := constant.InventoryTypeEquipment; invType <= constant.InventoryTypeCash; invType++ {
		items, ok := p.Tabs[invType]
		if ok == false {
			continue
		}
		writer.WriteU8(uint8(len(items)))
		for _, item := range items {
			item.Serialize(writer, dto.ItemSerializeOption{
				SlotMode: dto.SlotEncodeOmit,
			})
		}
	}
	return nil
}

func (p *Storage) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.StorageResult(reader.ReadU8())
	if p.Result == pconst.StorageResultOpen {
		p.NpcID = reader.ReadU32()
	}
	p.Slots = reader.ReadU8()
	mask := reader.ReadU64()
	if mask&2 != 0 {
		meso := reader.Read32()
		p.Meso = &meso
	}
	p.Tabs = map[constant.InventoryType][]dto.Item{}
	for invType := constant.InventoryTypeEquipment; invType <= constant.InventoryTypeCash; invType++ {
		if mask&(2<<invType) == 0 {
			continue
		}
		count := int(reader.ReadU8())
		items := make([]dto.Item, 0, count)
		for i := 0; i < count; i++ {
			items = append(items, dto.NewItemFromStream(reader))
		}
		p.Tabs[invType] = items
	}
}

type StorageError struct {
	Result pconst.StorageResult
}

func (p *StorageError) Opcode() uint16 {
	return 0xEA
}

func (p *StorageError) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(p.Result))
	return nil
}

func (p *StorageError) Deserialize(reader *stream.StreamReader) {
	p.Result = pconst.StorageResult(reader.ReadU8())
}

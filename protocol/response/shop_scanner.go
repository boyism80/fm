package response

import (
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type ShopScannerPopular struct {
	ItemIDs []uint32
}

func (p *ShopScannerPopular) Opcode() uint16 {
	return 0x35
}

func (p *ShopScannerPopular) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.ShopScannerResultPopular))
	writer.WriteU8(uint8(len(p.ItemIDs)))
	for _, itemID := range p.ItemIDs {
		writer.WriteU32(itemID)
	}
	return nil
}

func (p *ShopScannerPopular) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	count := int(reader.ReadU8())
	p.ItemIDs = make([]uint32, 0, count)
	for i := 0; i < count; i++ {
		p.ItemIDs = append(p.ItemIDs, reader.ReadU32())
	}
}

type ShopScannerEntry struct {
	OwnerName     string
	MapID         uint32
	Title         string
	PerBundle     uint32
	Bundles       uint32
	Price         int32
	SN            uint32
	Channel       int8
	InventoryType uint8
	Item          dto.Item
}

type ShopScannerResult struct {
	ItemID  uint32
	Entries []ShopScannerEntry
}

func (p *ShopScannerResult) Opcode() uint16 {
	return 0x35
}

func (p *ShopScannerResult) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(uint8(pconst.ShopScannerResultSearch))
	writer.WriteU32(0)
	writer.WriteU32(p.ItemID)
	writer.WriteU32(uint32(len(p.Entries)))
	for _, entry := range p.Entries {
		writer.WriteStr16(entry.OwnerName)
		writer.WriteU32(entry.MapID)
		writer.WriteStr16(entry.Title)
		writer.WriteU32(entry.PerBundle)
		writer.WriteU32(entry.Bundles)
		writer.Write32(entry.Price)
		writer.WriteU32(entry.SN)
		writer.Write8(entry.Channel)
		writer.WriteU8(entry.InventoryType)
		if entry.Item != nil {
			entry.Item.Serialize(writer, dto.ItemSerializeOption{SlotMode: dto.SlotEncodeOmit})
		}
	}
	return nil
}

func (p *ShopScannerResult) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	reader.ReadU32()
	p.ItemID = reader.ReadU32()
	count := int(reader.ReadU32())
	p.Entries = make([]ShopScannerEntry, 0, count)
	for i := 0; i < count; i++ {
		entry := ShopScannerEntry{
			OwnerName:     reader.ReadStr16(),
			MapID:         reader.ReadU32(),
			Title:         reader.ReadStr16(),
			PerBundle:     reader.ReadU32(),
			Bundles:       reader.ReadU32(),
			Price:         reader.Read32(),
			SN:            reader.ReadU32(),
			Channel:       reader.Read8(),
			InventoryType: reader.ReadU8(),
		}
		if entry.InventoryType == 1 {
			entry.Item = dto.NewItemFromStream(reader)
		}
		p.Entries = append(p.Entries, entry)
	}
}

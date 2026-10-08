package request

import "github.com/boyism80/fm/stream"

type ShopScannerOpen struct {
	Mode uint8
}

func (*ShopScannerOpen) Opcode() byte { return 0x31 }

func (p *ShopScannerOpen) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU8(p.Mode)
	return nil
}

func (p *ShopScannerOpen) Deserialize(reader *stream.StreamReader) {
	p.Mode = reader.ReadU8()
}

type ShopScannerWarp struct {
	SN    uint32
	MapID uint32
}

func (*ShopScannerWarp) Opcode() byte { return 0x32 }

func (p *ShopScannerWarp) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.SN)
	writer.WriteU32(p.MapID)
	return nil
}

func (p *ShopScannerWarp) Deserialize(reader *stream.StreamReader) {
	p.SN = reader.ReadU32()
	p.MapID = reader.ReadU32()
}

type UseShopScanner struct {
	Slot      int16
	ItemID    uint32
	SearchID  uint32
	HighFirst bool
}

func (*UseShopScanner) Opcode() byte { return 0x42 }

func (p *UseShopScanner) Serialize(writer *stream.StreamWriter) error {
	writer.Write16(p.Slot)
	writer.WriteU32(p.ItemID)
	writer.WriteU32(p.SearchID)
	writer.WriteBoolean(p.HighFirst)
	writer.WriteU32(0)
	return nil
}

func (p *UseShopScanner) Deserialize(reader *stream.StreamReader) {
	p.Slot = reader.Read16()
	p.ItemID = reader.ReadU32()
	p.SearchID = reader.ReadU32()
	p.HighFirst = reader.ReadU8() != 0
}

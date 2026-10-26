package response

import "github.com/boyism80/fm/stream"

// IncubatorResult shows what a pigmy egg hatched into; the client shows the failure popup when ItemID is 0.
type IncubatorResult struct {
	ItemID uint32
	Count  int16
}

func (p *IncubatorResult) Opcode() uint16 {
	return 0x34
}

func (p *IncubatorResult) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU32(p.ItemID)
	sw.Write16(p.Count)
	return nil
}

func (p *IncubatorResult) Deserialize(reader *stream.StreamReader) {
	p.ItemID = reader.ReadU32()
	p.Count = reader.Read16()
}

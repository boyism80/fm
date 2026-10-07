package request

import "github.com/boyism80/fm/stream"

type TeleportStoneTarget struct {
	MapID uint32
	Name  string
}

func (p *TeleportStoneTarget) Serialize(writer *stream.StreamWriter) {
	if p.Name != "" {
		writer.WriteU8(1)
		writer.WriteStr16(p.Name)
	} else {
		writer.WriteU8(0)
		writer.WriteU32(p.MapID)
	}
	writer.WriteU32(0)
}

func (p *TeleportStoneTarget) Deserialize(reader *stream.StreamReader) {
	if reader.Remaining() <= 4 {
		return
	}
	if reader.ReadU8() == 1 {
		p.Name = reader.ReadStr16()
	} else {
		p.MapID = reader.ReadU32()
	}
}

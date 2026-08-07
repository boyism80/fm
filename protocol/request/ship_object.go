package request

import "github.com/boyism80/fm/stream"

type ShipObject struct {
	MapID uint32
}

func (*ShipObject) Opcode() byte {
	return 0xB3
}

func (p *ShipObject) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(p.MapID)
	return nil
}

func (p *ShipObject) Deserialize(reader *stream.StreamReader) {
	p.MapID = reader.ReadU32()
}

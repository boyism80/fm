package response

import "github.com/boyism80/fm/stream"

type Tremble struct {
	Type  uint8
	Delay int32
}

func (p *Tremble) Opcode() uint16 {
	return 0x5F
}

func (p *Tremble) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(uint8(EnvironmentChangeModeTremble))
	sw.WriteU8(p.Type)
	sw.Write32(p.Delay)
	return nil
}

func (p *Tremble) Deserialize(reader *stream.StreamReader) {
	reader.ReadU8()
	p.Type = reader.ReadU8()
	p.Delay = reader.Read32()
}
